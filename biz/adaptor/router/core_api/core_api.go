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
	{
		_auth := root.Group("/auth", _authMw()...)
		_auth.GET("/checkEmail", append(_checkemailMw(), core_api.CheckEmail)...)
		_auth.POST("/emailLogin", append(_emailloginMw(), core_api.EmailLogin)...)
		_auth.GET("/giteeLogin", append(_giteeloginMw(), core_api.GiteeLogin)...)
		_auth.GET("/githubLogin", append(_githubloginMw(), core_api.GithubLogin)...)
		_auth.POST("/refreshToken", append(_refreshtokenMw(), core_api.RefreshToken)...)
		_auth.POST("/register", append(_registerMw(), core_api.Register)...)
		_auth.POST("/sendEmail", append(_sendemailMw(), core_api.SendEmail)...)
		_auth.POST("/setPasswordByEmail", append(_setpasswordbyemailMw(), core_api.SetPasswordByEmail)...)
		_auth.POST("/setPasswordByPassword", append(_setpasswordbypasswordMw(), core_api.SetPasswordByPassword)...)
	}
	{
		_content := root.Group("/content", _contentMw()...)
		_content.POST("/addFileToPublicSpace", append(_addfiletopublicspaceMw(), core_api.AddFileToPublicSpace)...)
		_content.POST("/askDownloadFile", append(_askdownloadfileMw(), core_api.AskDownloadFile)...)
		_content.POST("/askUploadAvatar", append(_askuploadavatarMw(), core_api.AskUploadAvatar)...)
		_content.POST("/askUploadFile", append(_askuploadfileMw(), core_api.AskUploadFile)...)
		_content.POST("/completelyRemoveFile", append(_completelyremovefileMw(), core_api.CompletelyRemoveFile)...)
		_content.POST("/createFile", append(_createfileMw(), core_api.CreateFile)...)
		_content.POST("/createPost", append(_createpostMw(), core_api.CreatePost)...)
		_content.POST("/createShareCode", append(_createsharecodeMw(), core_api.CreateShareCode)...)
		_content.POST("/createZone", append(_createzoneMw(), core_api.CreateZone)...)
		_content.POST("/deleteFile", append(_deletefileMw(), core_api.DeleteFile)...)
		_content.POST("/deletePost", append(_deletepostMw(), core_api.DeletePost)...)
		_content.POST("/deleteShareCode", append(_deletesharecodeMw(), core_api.DeleteShareCode)...)
		_content.POST("/deleteZone", append(_deletezoneMw(), core_api.DeleteZone)...)
		_content.GET("/getFileBySharingCode", append(_getfilebysharingcodeMw(), core_api.GetFileBySharingCode)...)
		_content.GET("/getOtherPost", append(_getotherpostMw(), core_api.GetOtherPost)...)
		_content.GET("/getOtherPosts", append(_getotherpostsMw(), core_api.GetOtherPosts)...)
		_content.GET("/getOwnPost", append(_getownpostMw(), core_api.GetOwnPost)...)
		_content.GET("/getOwnPosts", append(_getownpostsMw(), core_api.GetOwnPosts)...)
		_content.GET("/getPrivateFile", append(_getprivatefileMw(), core_api.GetPrivateFile)...)
		_content.GET("/getPrivateFiles", append(_getprivatefilesMw(), core_api.GetPrivateFiles)...)
		_content.GET("/getPublicFile", append(_getpublicfileMw(), core_api.GetPublicFile)...)
		_content.GET("/getPublicFiles", append(_getpublicfilesMw(), core_api.GetPublicFiles)...)
		_content.GET("/getRecycleBinFiles", append(_getrecyclebinfilesMw(), core_api.GetRecycleBinFiles)...)
		_content.GET("/getShareList", append(_getsharelistMw(), core_api.GetShareList)...)
		_content.GET("/getUser", append(_getuserMw(), core_api.GetUser)...)
		_content.GET("/getUserDetail", append(_getuserdetailMw(), core_api.GetUserDetail)...)
		_content.GET("/getZone", append(_getzoneMw(), core_api.GetZone)...)
		_content.POST("/moveFile", append(_movefileMw(), core_api.MoveFile)...)
		_content.GET("/parsingShareCode", append(_parsingsharecodeMw(), core_api.ParsingShareCode)...)
		_content.POST("/recoverRecycleBinFile", append(_recoverrecyclebinfileMw(), core_api.RecoverRecycleBinFile)...)
		_content.POST("/saveFileToPrivateSpace", append(_savefiletoprivatespaceMw(), core_api.SaveFileToPrivateSpace)...)
		_content.GET("/searchUser", append(_searchuserMw(), core_api.SearchUser)...)
		_content.POST("/updateFile", append(_updatefileMw(), core_api.UpdateFile)...)
		_content.POST("/updatePost", append(_updatepostMw(), core_api.UpdatePost)...)
		_content.POST("/updateUser", append(_updateuserMw(), core_api.UpdateUser)...)
		_content.POST("/updateZone", append(_updatezoneMw(), core_api.UpdateZone)...)
	}
	{
		_relation := root.Group("/relation", _relationMw()...)
		_relation.POST("/createRelation", append(_createrelationMw(), core_api.CreateRelation)...)
		_relation.POST("/deleteRelation", append(_deleterelationMw(), core_api.DeleteRelation)...)
		_relation.GET("/getFromRelations", append(_getfromrelationsMw(), core_api.GetFromRelations)...)
		_relation.GET("/getRelation", append(_getrelationMw(), core_api.GetRelation)...)
		_relation.GET("/getToRelations", append(_gettorelationsMw(), core_api.GetToRelations)...)
	}
}
