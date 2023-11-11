package app

import (
	"strconv"

	"github.com/Terence1105/Tinder/cmd/tinder/app/dto"
	redisdto "github.com/Terence1105/Tinder/pkg/storage/redis/tinder/dto"
	"github.com/Terence1105/Tinder/pkg/types"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// @Summary 加入用戶
// @Description 加入用戶和配對
// @Tags tinder
// @Accept json
// @Produce json
// @Param Request body dto.AddSinglePersonAndMatchReq true "raw"
// @Success 200 {object} BaseResponse "ok"
// @Failure 400 {object} BaseResponse "bad request"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/add-single-person-and-match [post]
func (a *App) AddSinglePersonAndMatch(c *gin.Context) {
	req := &dto.AddSinglePersonAndMatchReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	dto := &redisdto.Person{
		Name:       req.Name,
		Height:     req.Height,
		Gender:     req.Gender,
		DateCounts: req.DateCounts,
	}

	err := a.storage.AddPerson(c, dto)
	if err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	c.JSON(200, BaseResponse{Result: "success"})
}

// @Summary 移除用戶
// @Description 移除用戶
// @Tags tinder
// @Accept json
// @Produce json
// @Param Request body dto.RemoveSinglePersonReq true "raw"
// @Success 200 {object} BaseResponse "ok"
// @Failure 400 {object} BaseResponse "bad request"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/remove-single-person [post]
func (a *App) RemoveSinglePerson(c *gin.Context) {
	req := &dto.RemoveSinglePersonReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	err := a.storage.RemovePerson(c, req.Name, req.Gender)
	if err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	c.JSON(200, BaseResponse{Result: "success"})
}

// @Summary 配對n組
// @Description 配對n組
// @Tags tinder
// @Accept json
// @Produce json
// @Param Request body dto.QuerySinglePeopleReq true "raw"
// @Success 200 {object} BaseResponse "ok"
// @Failure 400 {object} BaseResponse "bad request"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/query-single-people [get]
func (a *App) QuerySinglePeople(c *gin.Context) {
	req := &dto.QuerySinglePeopleReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	matches := make([]dto.Match, 0, req.Counts)

	boys, err := a.storage.GetPeople(c, types.MinHeight, types.MaxHeight, req.Counts, types.BOY)
	if err != nil {
		c.JSON(400, BaseResponse{Error: err.Error()})
		return
	}

	for _, b := range boys {
		r, err := a.storage.GetDateCount(c, b.Name)
		if err != nil {
			c.JSON(400, BaseResponse{Error: err.Error()})
			return
		}

		dateCount, ok := strconv.Atoi(r)
		if ok != nil {
			c.JSON(400, BaseResponse{Error: err.Error()})
			return
		}

		girls, err := a.storage.GetPeople(c, types.MinHeight, b.Height, dateCount, types.GIRL)
		if err != nil {
			c.JSON(400, BaseResponse{Error: err.Error()})
			return
		}

		if girls == nil {
			break
		}

		for _, girl := range girls {
			r, err := a.storage.DecrementDateCount(c, girl.Name)
			if err != nil {
				c.JSON(400, BaseResponse{Error: err.Error()})
				return
			}

			if r == 0 {
				err = a.storage.RemovePerson(c, girl.Name, girl.Gender)
				if err != nil {
					c.JSON(400, BaseResponse{Error: err.Error()})
					return
				}
			}

			bdc, err := a.storage.DecrementDateCount(c, b.Name)
			if err != nil {
				c.JSON(400, BaseResponse{Error: err.Error()})
				return
			}

			if bdc == 0 {
				err = a.storage.RemovePerson(c, b.Name, b.Gender)
				if err != nil {
					c.JSON(400, BaseResponse{Error: err.Error()})
					return
				}
			}

			matches = append(matches, dto.Match{
				Boy:  b.Name,
				Girl: girl.Name,
			})

			if len(matches) >= req.Counts {
				break
			}
		}
	}

	resp := &dto.QuerySinglePeopleResp{
		Matches: matches,
	}

	c.JSON(200, BaseResponse{Result: resp})
}
