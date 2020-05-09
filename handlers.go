package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func processGetPrivatePage(ctx *fasthttp.RequestCtx) {
	var userId UserId
	id, err := ctx.QueryArgs().GetUint("userId")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	userId.Id = int64(id)
	user, err := GetPrivatePage(userId.Id)
	if err != nil {
		writeError(ctx, "error with getting my page", http.StatusBadRequest)
		return
	}
	response, _ := json.Marshal(user)
	ctx.Write(response)
}

func processGetCommonPage(ctx *fasthttp.RequestCtx) {
	commonUser := CommonUser{}
	userId, err := ctx.QueryArgs().GetUint("userId")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	commonId, err := ctx.QueryArgs().GetUint("commonId")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	if userId == commonId {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	commonUser.Id = int64(commonId)
	commonUser, err = GetCommonPage(int64(userId), int64(commonId))
	if err != nil {
		writeError(ctx, "error with getting my page", http.StatusBadRequest)
		return
	}
	response, _ := json.Marshal(commonUser)
	ctx.Write(response)
}

func processLogin(ctx *fasthttp.RequestCtx) {
	var login Login
	regLogin := ctx.Request.Body()
	err := json.Unmarshal(regLogin, &login)
	if err != nil {
		log.Println(err)
		writeError(ctx, "error with login", http.StatusBadRequest)
		return
	}
	if ok := LoginPass(login); !ok {
		writeError(ctx, `{"isAuth":false}`, http.StatusForbidden)
		return
	}
	response := []byte(`{"isAuth":true}`)
	ctx.Write(response)
}

func processRegistration(ctx *fasthttp.RequestCtx) {
	var user User
	regUser := ctx.Request.Body()
	err := json.Unmarshal(regUser, &user)
	if err != nil {
		log.Println(err)
		writeError(ctx, "error with unmarshalling", http.StatusBadRequest)
		return
	}
	isUnique := IsUnique(user.Email)
	if !isUnique {
		writeError(ctx, "Your Email is not unique", http.StatusForbidden)
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	log.Println(string(hashedPassword))
	user.Id, err = Registration(user)
	if err != nil {
		writeError(ctx, "error with registration", http.StatusBadRequest)
		return
	}
	response, _ := json.Marshal(user.UserId)
	ctx.Write(response)
}

func processGetFriends(ctx *fasthttp.RequestCtx) {
	var userId UserId
	id, err := ctx.QueryArgs().GetUint("userId")
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	userId.Id = int64(id)
	friends, err := GetFriends(userId.Id)
	if err != nil {
		writeError(ctx, "error with getting friends", http.StatusBadRequest)
		return
	}
	var response []byte
	if len(friends) > 0 {
		response, _ = json.Marshal(friends)
	} else {
		response = []byte(`[]`)
	}
	ctx.Write(response)

}

func processAddFriend(ctx *fasthttp.RequestCtx) {
	var uf UserFriend
	regUf := ctx.Request.Body()
	err := json.Unmarshal(regUf, &uf)
	if err != nil {
		log.Println(err)
		writeError(ctx, "error with unmarshalling", http.StatusBadRequest)
		return
	}
	err = AddFriend(uf.UserId, uf.FriendId)
	if err != nil {
		writeError(ctx, "error with adding friend", http.StatusBadRequest)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}

func processRemoveFriend(ctx *fasthttp.RequestCtx) {
	var uf UserFriend
	regUf := ctx.Request.Body()
	err := json.Unmarshal(regUf, &uf)
	if err != nil {
		log.Println(err)
		writeError(ctx, "error with unmarshalling", http.StatusBadRequest)
		return
	}
	err = RemoveFriend(uf.UserId, uf.FriendId)
	if err != nil {
		writeError(ctx, "error with removing friend", http.StatusBadRequest)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
}