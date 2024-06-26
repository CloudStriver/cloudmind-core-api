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
		_auth.POST("/askUploadAvatar", append(_askuploadavatarMw(), core_api.AskUploadAvatar)...)
		_auth.GET("/checkEmail", append(_checkemailMw(), core_api.CheckEmail)...)
		_auth.POST("/emailLogin", append(_emailloginMw(), core_api.EmailLogin)...)
		_auth.GET("/qqLogin", append(_qqloginMw(), core_api.QQLogin)...)
		_auth.POST("/refreshToken", append(_refreshtokenMw(), core_api.RefreshToken)...)
		_auth.POST("/register", append(_registerMw(), core_api.Register)...)
		_auth.POST("/sendEmail", append(_sendemailMw(), core_api.SendEmail)...)
		_auth.POST("/setPasswordByEmail", append(_setpasswordbyemailMw(), core_api.SetPasswordByEmail)...)
		_auth.POST("/setPasswordByPassword", append(_setpasswordbypasswordMw(), core_api.SetPasswordByPassword)...)
		_auth.POST("/weixinCallback", append(_weixincallbackMw(), core_api.WeixinCallBack)...)
		_auth.POST("/weixinIsLogin", append(_weixinisloginMw(), core_api.WeixinIsLogin)...)
		_auth.GET("/weixinLogin", append(_weixinloginMw(), core_api.WeixinLogin)...)
	}
	{
		_comment := root.Group("/comment", _commentMw()...)
		_comment.POST("/createComment", append(_createcommentMw(), core_api.CreateComment)...)
		_comment.POST("/deleteComment", append(_deletecommentMw(), core_api.DeleteComment)...)
		_comment.POST("/deleteCommentSubject", append(_deletecommentsubjectMw(), core_api.DeleteCommentSubject)...)
		_comment.GET("/getComment", append(_getcommentMw(), core_api.GetComment)...)
		_comment.GET("/getCommentBlocks", append(_getcommentblocksMw(), core_api.GetCommentBlocks)...)
		_comment.GET("/getCommentSubject", append(_getcommentsubjectMw(), core_api.GetCommentSubject)...)
		_comment.GET("/getComments", append(_getcommentsMw(), core_api.GetComments)...)
		_comment.POST("/setCommentAttrs", append(_setcommentattrsMw(), core_api.SetCommentAttrs)...)
		_comment.POST("/updateComment", append(_updatecommentMw(), core_api.UpdateComment)...)
		_comment.POST("/updateCommentSubject", append(_updatecommentsubjectMw(), core_api.UpdateCommentSubject)...)
	}
	{
		_content := root.Group("/content", _contentMw()...)
		_content.POST("/addFileToPublicSpace", append(_addfiletopublicspaceMw(), core_api.AddFileToPublicSpace)...)
		_content.POST("/askDownloadFile", append(_askdownloadfileMw(), core_api.AskDownloadFile)...)
		_content.POST("/askUploadFile", append(_askuploadfileMw(), core_api.AskUploadFile)...)
		_content.POST("/checkFile", append(_checkfileMw(), core_api.CheckFile)...)
		_content.POST("/completelyRemoveFile", append(_completelyremovefileMw(), core_api.CompletelyRemoveFile)...)
		_content.POST("/createFeedBack", append(_createfeedbackMw(), core_api.CreateFeedBack)...)
		_content.POST("/createFile", append(_createfileMw(), core_api.CreateFile)...)
		_content.POST("/createPost", append(_createpostMw(), core_api.CreatePost)...)
		_content.POST("/createProduct", append(_createproductMw(), core_api.CreateProduct)...)
		_content.POST("/createShareCode", append(_createsharecodeMw(), core_api.CreateShareCode)...)
		_content.POST("/deleteFile", append(_deletefileMw(), core_api.DeleteFile)...)
		_content.POST("/deletePost", append(_deletepostMw(), core_api.DeletePost)...)
		_content.POST("/deleteProduct", append(_deleteproductMw(), core_api.DeleteProduct)...)
		_content.POST("/deleteShareCode", append(_deletesharecodeMw(), core_api.DeleteShareCode)...)
		_content.POST("/emptyRecycleBin", append(_emptyrecyclebinMw(), core_api.EmptyRecycleBin)...)
		_content.GET("/getFileBySharingCode", append(_getfilebysharingcodeMw(), core_api.GetFileBySharingCode)...)
		_content.GET("/getLatestRecommend", append(_getlatestrecommendMw(), core_api.GetLatestRecommend)...)
		_content.GET("/getPopularRecommend", append(_getpopularrecommendMw(), core_api.GetPopularRecommend)...)
		_content.GET("/getPost", append(_getpostMw(), core_api.GetPost)...)
		_content.GET("/getPosts", append(_getpostsMw(), core_api.GetPosts)...)
		_content.GET("/getPrivateFile", append(_getprivatefileMw(), core_api.GetPrivateFile)...)
		_content.GET("/getPrivateFiles", append(_getprivatefilesMw(), core_api.GetPrivateFiles)...)
		_content.GET("/getProduct", append(_getproductMw(), core_api.GetProduct)...)
		_content.GET("/getProducts", append(_getproductsMw(), core_api.GetProducts)...)
		_content.GET("/getPublicFile", append(_getpublicfileMw(), core_api.GetPublicFile)...)
		_content.GET("/getPublicFiles", append(_getpublicfilesMw(), core_api.GetPublicFiles)...)
		_content.GET("/getRecommendByItem", append(_getrecommendbyitemMw(), core_api.GetRecommendByItem)...)
		_content.GET("/getRecommendByUser", append(_getrecommendbyuserMw(), core_api.GetRecommendByUser)...)
		_content.GET("/getRecycleBinFiles", append(_getrecyclebinfilesMw(), core_api.GetRecycleBinFiles)...)
		_content.GET("/getShareList", append(_getsharelistMw(), core_api.GetShareList)...)
		_content.GET("/getUser", append(_getuserMw(), core_api.GetUser)...)
		_content.GET("/getUserDetail", append(_getuserdetailMw(), core_api.GetUserDetail)...)
		_content.POST("/makeFilePrivate", append(_makefileprivateMw(), core_api.MakeFilePrivate)...)
		_content.POST("/moveFile", append(_movefileMw(), core_api.MoveFile)...)
		_content.GET("/parsingShareCode", append(_parsingsharecodeMw(), core_api.ParsingShareCode)...)
		_content.POST("/recoverRecycleBinFile", append(_recoverrecyclebinfileMw(), core_api.RecoverRecycleBinFile)...)
		_content.POST("/saveFileToPrivateSpace", append(_savefiletoprivatespaceMw(), core_api.SaveFileToPrivateSpace)...)
		_content.GET("/searchUser", append(_searchuserMw(), core_api.SearchUser)...)
		_content.POST("/updateFile", append(_updatefileMw(), core_api.UpdateFile)...)
		_content.POST("/updatePost", append(_updatepostMw(), core_api.UpdatePost)...)
		_content.POST("/updateProduct", append(_updateproductMw(), core_api.UpdateProduct)...)
		_content.POST("/updateUser", append(_updateuserMw(), core_api.UpdateUser)...)
	}
	{
		_label := root.Group("/label", _labelMw()...)
		_label.POST("/createLabel", append(_createlabelMw(), core_api.CreateLabel)...)
		_label.POST("/deleteLabel", append(_deletelabelMw(), core_api.DeleteLabel)...)
		_label.GET("/getLabel", append(_getlabelMw(), core_api.GetLabel)...)
		_label.GET("/getLabels", append(_getlabelsMw(), core_api.GetLabels)...)
		_label.GET("/getLabelsInBatch", append(_getlabelsinbatchMw(), core_api.GetLabelsInBatch)...)
		_label.POST("/updateLabel", append(_updatelabelMw(), core_api.UpdateLabel)...)
	}
	{
		_rank := root.Group("/rank", _rankMw()...)
		_rank.GET("/getHotRanks", append(_gethotranksMw(), core_api.GetHotRanks)...)
	}
	{
		_relation := root.Group("/relation", _relationMw()...)
		_relation.POST("/createRelation", append(_createrelationMw(), core_api.CreateRelation)...)
		_relation.POST("/deleteRelation", append(_deleterelationMw(), core_api.DeleteRelation)...)
		_relation.GET("/getFromRelations", append(_getfromrelationsMw(), core_api.GetFromRelations)...)
		_relation.GET("/getRelation", append(_getrelationMw(), core_api.GetRelation)...)
		_relation.GET("/getRelationPaths", append(_getrelationpathsMw(), core_api.GetRelationPaths)...)
		_relation.GET("/getToRelations", append(_gettorelationsMw(), core_api.GetToRelations)...)
	}
	{
		_system := root.Group("/system", _systemMw()...)
		_system.POST("/createSlider", append(_createsliderMw(), core_api.CreateSlider)...)
		_system.POST("/deleteNotifications", append(_deletenotificationsMw(), core_api.DeleteNotifications)...)
		_system.POST("/deleteSlider", append(_deletesliderMw(), core_api.DeleteSlider)...)
		_system.GET("/getNotificationCount", append(_getnotificationcountMw(), core_api.GetNotificationCount)...)
		_system.GET("/getNotifications", append(_getnotificationsMw(), core_api.GetNotifications)...)
		_system.GET("/getSliders", append(_getslidersMw(), core_api.GetSliders)...)
		_system.POST("/updateSlider", append(_updatesliderMw(), core_api.UpdateSlider)...)
	}
}
