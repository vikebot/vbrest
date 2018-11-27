package vbapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/vikebot/vbcore"
	"github.com/vikebot/vbdb"
	"github.com/vikebot/vbnet"
	"github.com/vikebot/vbrest/vbmail"
	"go.uber.org/zap"
)

var (
	registercodeValidator = regexp.MustCompile("^[a-zA-Z0-9_-]{32}$")
)

type RegisterConfirmRequest struct {
	Code         *string      `json:"code"`
	User         *vbcore.User `json:"user"`
	Verification *string      `json:"verification"`
	Recaptcha    *string      `json:"recaptcha"`
}

// RegisterConfirm registers the user specified by the user object
func RegisterConfirm(data RegisterConfirmRequest, ip string, ctx *zap.Logger) error {
	// Check recaptcha
	if data.Recaptcha == nil {
		return vbnet.NewHTTPError("Recaptcha not ticked", http.StatusBadRequest, codeRecaptchaNotTicked, nil)
	}
	hasTicked, err := recaptcha.Confirm(ip, *data.Recaptcha)
	if err != nil {
		return errInternalServerError
	}
	if !hasTicked {
		return vbnet.NewHTTPError("Recaptcha not ticked", http.StatusBadRequest, codeRecaptchaNotTicked, nil)
	}

	// Check registration code syntax
	if data.Code == nil || !registercodeValidator.MatchString(*data.Code) {
		return vbnet.NewHTTPError("Code must be valid", http.StatusBadRequest, codeInvalidRegisterCode, nil)
	}

	// Verify that user is provided
	if data.User == nil {
		return vbnet.NewHTTPError("User cannot be null", http.StatusBadRequest, codeUserCannotBeNull, nil)
	}

	// Validate the user provided object
	user, valid := data.User.Validate()
	if !valid {
		return vbnet.NewHTTPError("User state is invalid", http.StatusBadRequest, codeBadUserState, nil)
	}

	// Load id of originial user from provided reg code
	userID, finished, success := vbdb.UserIDFromRegcodeCtx(*data.Code, ctx)
	if !success {
		return errInternalServerError
	}
	if finished {
		return vbnet.NewHTTPError("You already finished registration", http.StatusBadRequest, codeAlreadyFinishedRegistration, nil)
	}
	if userID == 0 {
		return vbnet.NewHTTPError("Registration code unknown", http.StatusBadRequest, codeRegistrationCodeUnknown, nil)
	}

	// Load olduser with id
	oldUser, success := vbdb.UserFromIDCtx(userID, ctx)
	if !success {
		return errInternalServerError
	}

	// Verify that email addresses aren't "new"
	var selected bool
	var selectedPrimary vbcore.Email
	for idx, ne := range user.Emails {
		if ne.Email != oldUser.Emails[idx].Email {
			return vbnet.NewHTTPError("Email addresses manipulated", http.StatusBadRequest, codeEmailAddressManipulated, nil)
		}
		if oldUser.Emails[idx].Status == 0 && ne.Status == 1 {
			return vbnet.NewHTTPError("Email address status manipulated", http.StatusBadRequest, codeEmailStatusManipulated, nil)
		}
		if ne.Status == vbcore.EmailPrimary {
			if !selected {
				selected = true
				selectedPrimary = ne
			} else {
				return vbnet.NewHTTPError("Cannot use multiple primary addresses", http.StatusBadRequest, codeCannotUseMultiplePrimaryEmail, nil)
			}
		}
	}

	// Check if primary was selected
	if !selected {
		return vbnet.NewHTTPError("Must have primary email. Request manipulated", http.StatusBadRequest, codeMustHavePrimaryEmail, nil)
	}

	// If user has selected a primary address and it's verified already
	if selectedPrimary.Status == vbcore.EmailVerified {
		success = vbdb.UpdateUserEmailStatusCtx(userID, selectedPrimary.Email, vbcore.EmailPrimary, ctx)
		if !success {
			return errInternalServerError
		}
	} else {
		if data.Verification == nil {
			last, valid, success := vbdb.UserEmailVerificationLoadCtx(userID, selectedPrimary.Email, ctx)
			if !success {
				return errInternalServerError
			}

			if valid {
				sec := time.Now().UTC().Sub(*last).Seconds()
				if sec < 60*5 {
					return vbnet.NewHTTPError("Unable to send verification email. Quota for user exhausted. Try again in "+strconv.Itoa(int(sec))+" seconds.", http.StatusBadRequest, codeEmailQuotaExhausted, nil)
				}
			}

			verificationCode, err := vbcore.CryptoGenString(8)
			if err != nil {
				ctx.Error("", zap.Error(err))
				return errInternalServerError
			}
			verificationCode = strings.ToLower(verificationCode[:4] + " " + verificationCode[4:])

			success = vbdb.UserEmailVerificationSetCtx(userID, selectedPrimary.Email, verificationCode, ctx)
			if !success {
				return errInternalServerError
			}

			plainText := fmt.Sprintf("Dear %s,\nIn order to verify this email address with Vikebot (https://vikebot.com) use the verification code: %s\nYour Vikebot Team!\n\n\nIf you didn't register with us, you can ignore this email.", user.Name, verificationCode)
			htmlText := fmt.Sprintf("<h3>Dear %s,</h3><p>In order to verify this email address with Vikebot (https://vikebot.com) use the verification code: <strong>%s</strong></p><p>Your Vikebot Team</p><br><br><p>If you didn't register with us, you can ignore this email.</p>", user.Name, verificationCode)

			err = vbmail.SendTo("[Action Required] Verify your Email with Vikebot", user.Name, selectedPrimary.Email, plainText, htmlText)
			if err != nil {
				ctx.Error("Unable to send email",
					zap.Error(err),
					zap.String("receiver", selectedPrimary.Email))
				return errInternalServerError
			}

			return vbnet.NewHTTPError("Please enter the code sent to your primary email address in the verification box.", http.StatusExpectationFailed, codeRegisterVerificationEntry, nil)
		} else {
			verified, success := vbdb.UserEmailVerificationIsCtx(userID, selectedPrimary.Email, *data.Verification, ctx)
			if !success {
				return errInternalServerError
			}
			if !verified {
				return vbnet.NewHTTPError("Invalid verification code", http.StatusBadRequest, codeInvalidEmailVerificationCode, nil)
			}

			success = vbdb.UpdateUserEmailStatusCtx(userID, selectedPrimary.Email, vbcore.EmailPrimary, ctx)
			if !success {
				return errInternalServerError
			}
		}
	}

	// Verify that web addresses havn't changed
	for _, item := range user.Web {
		var found bool
		for _, old := range oldUser.Web {
			if item == old {
				found = true
				break
			}
		}
		if !found {
			return vbnet.NewHTTPError("Invalid web link. Manipulated", http.StatusBadRequest, codeManipulatedWebLink, nil)
		}
	}
	success = vbdb.UserDeleteWebExpectCtx(userID, user.Web, ctx)
	if !success {
		return errInternalServerError
	}

	// Verify that social links havn't changed
	socialKeys := []string{}
	for k, v := range user.Social {
		if link, ok := oldUser.Social[k]; !ok || link != v {
			return vbnet.NewHTTPError("Invalid social platform. Request manipulated", http.StatusBadRequest, codeInvalidSocialPlatfrom, nil)
		}
		socialKeys = append(socialKeys, k)
	}
	success = vbdb.UserDeleteSocialExpectCtx(userID, socialKeys, ctx)
	if !success {
		return errInternalServerError
	}

	// Update auto prop fields
	success = vbdb.UpdateUserCtx(user.User(), oldUser, "register_autoPropUpdate", ctx)
	if !success {
		return errInternalServerError
	}

	// Set registration to done
	success = vbdb.UserSetRegistrationDoneCtx(userID, ctx)
	if !success {
		return errInternalServerError
	}

	return nil
}
