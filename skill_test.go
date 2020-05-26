package golexa_test

import (
	"context"
	"testing"

	"github.com/robsignorelli/golexa"
	"github.com/stretchr/testify/suite"
)

func TestSkillSuite(t *testing.T) {
	suite.Run(t, new(SkillSuite))
}

type SkillSuite struct {
	suite.Suite
}

func (suite SkillSuite) TestEmptyRouteTable() {
	skill := golexa.Skill{}

	_, err := skill.Handle(context.TODO(), golexa.NewIntentRequest("", golexa.NewSlots()))
	suite.Error(err, "Should result in an error when there's no intent name")

	_, err = skill.Handle(context.TODO(), golexa.NewIntentRequest("NotFound", golexa.NewSlots()))
	suite.Error(err, "Should result in an error when there's no route for the intent name")
}

func (suite SkillSuite) TestStandardRouteTable() {
	skill := golexa.Skill{}
	skill.RouteIntent("Intent1", func(ctx context.Context, request golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(request).Speak("Handler 1").Ok()
	})
	skill.RouteIntent("Intent2", func(ctx context.Context, request golexa.Request) (golexa.Response, error) {
		return golexa.NewResponse(request).Speak("Handler 2").Ok()
	})

	_, err := skill.Handle(context.TODO(), golexa.NewIntentRequest("NotFound", golexa.NewSlots()))
	suite.Error(err, "Should result in an error when there's no route for the intent name")

	res, err := skill.Handle(context.TODO(), golexa.NewIntentRequest("Intent1", golexa.NewSlots()))
	suite.NoError(err, "Should not generate error for valid routes")
	suite.Equal("<speak>Handler 1</speak>", res.Body.OutputSpeech.SSML,
		"Should execute the appropriate handler for valid intent names")

	res, err = skill.Handle(context.TODO(), golexa.NewIntentRequest("Intent2", golexa.NewSlots()))
	suite.NoError(err, "Should not generate error for valid routes")
	suite.Equal("<speak>Handler 2</speak>", res.Body.OutputSpeech.SSML,
		"Should execute the appropriate handler for valid intent names")
}
