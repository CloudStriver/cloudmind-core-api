// Code generated by hertz generator. DO NOT EDIT.

package core_api

import (
	core_api "github.com/CloudStriver/cloudmind-core-api/biz/adaptor/controller/core_api"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	root.POST("/sts", append(_applysignedurlMw(), core_api.ApplySignedUrl)...)
	{
		_auth := root.Group("/auth", _authMw()...)
		_auth.POST("/refresh", append(_refreshtokenMw(), core_api.RefreshToken)...)
		_auth.POST("/register", append(_registerMw(), core_api.Register)...)
		_auth.POST("/send", append(_sendemailMw(), core_api.SendEmail)...)
		{
			_login := _auth.Group("/login", _loginMw()...)
			_login.POST("/email", append(_emailloginMw(), core_api.EmailLogin)...)
			_login.GET("/gitee", append(_giteeloginMw(), core_api.GiteeLogin)...)
			_login.GET("/github", append(_githubloginMw(), core_api.GithubLogin)...)
		}
		{
			_reset := _auth.Group("/reset", _resetMw()...)
			_reset.POST("/email", append(_setpasswordbyemailMw(), core_api.SetPasswordByEmail)...)
			_reset.POST("/password", append(_setpasswordbypasswordMw(), core_api.SetPasswordByPassword)...)
		}
	}
	{
		_content := root.Group("/content", _contentMw()...)
		{
			_user := _content.Group("/user", _userMw()...)
			_user.GET("/search", append(_searchuserMw(), core_api.SearchUser)...)
			_user.POST("/update", append(_updateuserMw(), core_api.UpdateUser)...)
		}
	}
	root.GET("/relation", append(_getrelationMw(), core_api.GetRelation)...)
	_relation := root.Group("/relation", _relationMw()...)
	_relation.POST("/delete", append(_deleterelationMw(), core_api.DeleteRelation)...)
	{
		_relation0 := root.Group("/relation", _relation0Mw()...)
		_relation0.POST("/create", append(_createrelationMw(), core_api.CreateRelation)...)
	}
	{
		_relations := root.Group("/relations", _relationsMw()...)
		_relations.GET("/from", append(_getfromrelationsMw(), core_api.GetFromRelations)...)
		_relations.GET("/to", append(_gettorelationsMw(), core_api.GetToRelations)...)
	}
}
