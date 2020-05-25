package golexa_test

import (
	"encoding/json"
	"testing"

	"github.com/robsignorelli/golexa"
	"github.com/stretchr/testify/suite"
)

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

type RequestSuite struct {
	suite.Suite
}

func (suite RequestSuite) parseJSON(input string) golexa.Request {
	req := golexa.Request{}
	err := json.Unmarshal([]byte(input), &req)
	suite.Require().Nil(err)
	return req
}

func (suite RequestSuite) TestJSON_IdentityInfo() {
	var input = `{
		"version": "1.0",
		"session": {
			"new": true,
			"sessionId": "session.123",
			"application": {
				"applicationId": "skill.456"
			},
			"user": {
				"userId": "account.789",
				"accessToken": "token.acb"
			}
		},
		"context": {
			"System": {
				"application": {
					"applicationId": "skill.456"
				},
				"user": {
					"userId": "account.789",
					"accessToken": "token.acb"
				},
				"device": {
					"deviceId": "device.def",
					"supportedInterfaces": {}
				},
				"apiEndpoint": "https://api.amazonalexa.com",
				"apiAccessToken": "3yJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6IjEifQ.eyJhdWQiOiJodHRwczovL2FwaS5hbWF6b25hbGV4YS5jb20iLCJpc3MiOiJBbGV4YVNraWxsS2l0Iiwic3ViIjoiYW16bjEuYXNrLnNraWxsLmY0YWRlM2NmLWEzYmYtNDYyNS05M2ExLTgwYTM1MjNmMWQyYyIsImV4cCI6MTU1Mjc2NTg5OCwiaWF0IjoxNTUyNzY1NTk4LCJuYmYiOjE1NTI3NjU1OTgsInByaXZhdGVDbGFpbXMiOnsiY29uc2VudFRva2VuIjpudWxsLCJkZXZpY2VJZCI6ImFtem4xLmFzay5kZXZpY2UuQUg3V1dXRFJHVElVNkIyV1VHRklQUVNVQVNWWENUV01DTjNUTERIT00zUE81TzZYWUVDUFBEUFlUQTZWREU3U0tSTktHQUdWSEFQTU5PR0ZGSzVDM1RHR1VEV0lEWEdIWUFZREhRTEdLN05TTEZIN0Q3V1hUUjNHR1lTNlNaRTU0UUxITFBTSUNNSTM3TUJFRzRTTzMyUFlEVVNCSEEzQTRUNjNFTURERjJIQ1Y0WE1RWjU3QSIsInVzZXJJZCI6ImFtem4xLmFzay5hY2NvdW50LkFIQjRYRTNaQU9aMkFQVVlXUFlDUlJWSVhHQVpHNFZMNVJZVEEzTVBVNTJERDRDVlA0R1dTVFRHT0E3UldGQVk0RjZNUkZQWlVSNkhLQ1pFTUJMRkk0NVlRRVFBNk9KQUxSSE1OUU5LSEdHSFZMSFpRVEdCSlNRTTZXS0pCMk81UFhVWlpLNlYyQkczMk1CTjdRVlFYUzJRUjNSNllYRjMyRFFOR1pFVkpWTlNSVVpESDZXQkZONlBLSE5ISkJPTUdZNDZLNlhXUE1PVE02USJ9fQ.jNi3S0fY_txXfhiMPiB8X4IvVbHlM8mWOa1ao2fUn--OU-kPou2Lio5TMxDeUnbe7n191HdQ6n6NSU8Zepnqy1syi9al-F_uhjsG8MdadO5BOovIyxul3zBoKWrWqHXYeYZ2OSFZEihSNyy8IrweL4PzJ7WWSCcoFJXY9jcIhjszbSuWn_Ov354bz0WHPWM_Wx9qrJnSqmEvttCDv6xxRrCLrEehWqoVEspODWXgwL-FvPzUsjZdYeO17n7UbWQOIr4OokBSuZz-I7sS9Hn-KmjJ2bLieGZF5rRBTqzQPvGFMh0h5l0BipMn0wGgOZnFmM6GSqumw08weH8Iq4WwIZ"
			}
		},
		"request": {
			"type": "IntentRequest",
			"requestId": "request.890",
			"timestamp": "2019-03-16T19:46:38Z",
			"locale": "en-US",
			"intent": {
				"name": "FooIntent",
				"confirmationStatus": "NONE",
				"slots": {
					"name": {
						"name": "name",
						"value": "Bob Loblaw",
						"confirmationStatus": "NONE",
						"source": "USER"
					}
				}
			},
			"dialogState": "COMPLETED"
		}
	}`
	req := suite.parseJSON(input)
	suite.Equal("session.123", req.Session.ID,
		"Should populate session id properly")
	suite.Equal("device.def", req.Context.System.Device.ID,
		"Should populate device id properly")

	suite.Equal("skill.456", req.Session.Application.ID,
		"Should populate skill id properly (from Session)")
	suite.Equal("skill.456", req.Context.System.Application.ID,
		"Should populate skill id properly (from Context)")

	suite.Equal("account.789", req.Session.User.ID,
		"Should populate user info properly (from Session)")
	suite.Equal("token.acb", req.Session.User.AccessToken,
		"Should populate user info properly (from Session)")
	suite.Equal("account.789", req.Context.System.User.ID,
		"Should populate user info properly (from Context)")
	suite.Equal("token.acb", req.Context.System.User.AccessToken,
		"Should populate user info properly (from Context)")

	suite.Require().NotNil(req.Body.Intent,
		"Should populate non-nil intent body")
	suite.Equal("FooIntent", req.Body.Intent.Name,
		"Should populate intent name properly")

	// Make sure all of our shorthand helpers work
	suite.Equal("session.123", req.SessionID(),
		"Shorthand helper should fetch session id properly")
	suite.Equal("account.789", req.UserID(),
		"Shorthand helper should fetch user id properly")
	suite.Equal("token.acb", req.UserAccessToken(),
		"Shorthand helper should fetch access token properly")
	suite.Equal("device.def", req.DeviceID(),
		"Shorthand helper should fetch device id properly")
	suite.Equal("skill.456", req.SkillID(),
		"Shorthand helper should fetch skill id properly")
}

func (suite RequestSuite) TestJSON_NoSlots() {
	var input = `{
		"version": "1.0",
		"request": {
			"type": "IntentRequest",
			"requestId": "request.890",
			"timestamp": "2019-03-16T19:46:38Z",
			"locale": "en-US",
			"intent": {
				"name": "FooIntent",
				"confirmationStatus": "NONE",
				"slots": {}
			},
			"dialogState": "COMPLETED"
		}
	}`
	req := suite.parseJSON(input)
	suite.Require().NotNil(req.Body.Intent,
		"Should have an Intent body regardless of slots")
	suite.Len(req.Body.Intent.Slots, 0,
		"Should still have an empty slot map when there are no slots")
	suite.Equal("", req.Body.Intent.Slots["askldfjaslkdfj"].Value,
		"Should have have blank slot values for non-existent slots")
}

func (suite RequestSuite) TestJSON_OneSlot() {
	var input = `{
		"version": "1.0",
		"request": {
			"type": "IntentRequest",
			"requestId": "request.890",
			"timestamp": "2019-03-16T19:46:38Z",
			"locale": "en-US",
			"intent": {
				"name": "FooIntent",
				"confirmationStatus": "NONE",
				"slots": {
					"name": {
						"name": "name",
						"value": "Bob Loblaw",
						"confirmationStatus": "NONE",
						"source": "USER"
					}
				}
			},
			"dialogState": "COMPLETED"
		}
	}`
	req := suite.parseJSON(input)
	suite.Require().NotNil(req.Body.Intent,
		"Should have an Intent body regardless of slots")
	suite.Require().Len(req.Body.Intent.Slots, 1,
		"Should still have one slot in the map when the JSON contains 1 slot")
	suite.Equal("Bob Loblaw", req.Body.Intent.Slots["name"].Value,
		"Should have the proper slot value from the JSON")
	suite.Equal("", req.Body.Intent.Slots["askldfjaslkdfj"].Value,
		"Should have have blank slot values for non-existent slots")
}

func (suite RequestSuite) TestJSON_MultipleSlots() {
	var input = `{
		"version": "1.0",
		"request": {
			"type": "IntentRequest",
			"requestId": "request.890",
			"timestamp": "2019-03-16T19:46:38Z",
			"locale": "en-US",
			"intent": {
				"name": "FooIntent",
				"confirmationStatus": "NONE",
				"slots": {
					"name": {
						"name": "name",
						"value": "Bob Loblaw",
						"confirmationStatus": "NONE",
						"source": "USER"
					},
					"age": {
						"name": "age",
						"value": "",
						"confirmationStatus": "NONE",
						"source": "USER"
					},
					"hobby": {
						"name": "hobby",
						"value": "gaming",
						"resolutions": {
							"resolutionsPerAuthority": [
								{
									"authority": "authority.99",
									"status": {
										"code": "ER_SUCCESS_MATCH"
									},
									"values": [
										{
											"value": {
												"name": "video games",
												"id": "foobarbaz"
											}
										}
									]
								}
							]
						}
					}
				}
			},
			"dialogState": "COMPLETED"
		}
	}`
	req := suite.parseJSON(input)
	suite.Require().NotNil(req.Body.Intent,
		"Should have an Intent body regardless of slots")
	suite.Require().Len(req.Body.Intent.Slots, 3,
		"Should still have 3 slots in the map when the JSON contains 3 slot")

	suite.Equal("Bob Loblaw", req.Body.Intent.Slots["name"].Value,
		"Should have the proper slot value from the JSON")
	suite.Equal("Bob Loblaw", req.Body.Intent.Slots.Resolve("name"),
		"Should have resolve the proper slot value for slots in the request")

	suite.Equal("", req.Body.Intent.Slots["age"].Value,
		"Should have the proper slot value from the JSON")
	suite.Equal("", req.Body.Intent.Slots.Resolve("age"),
		"Should have resolve the proper slot value for slots in the request")

	suite.Equal("gaming", req.Body.Intent.Slots["hobby"].Value,
		"Should have the proper slot value from the JSON")
	suite.Equal("video games", req.Body.Intent.Slots.Resolve("hobby"),
		"Should have resolve the proper slot value for slots in the request")
}
