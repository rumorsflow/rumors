//go:build nop

package sys

//	@title		Rumors System API
//	@version	1.0
//	@BasePath	/sys/api

//	@securityDefinitions.apikey	SysAuth
//	@in							header
//	@name						Authorization

//	@Summary		Sign In
//	@Description	get sign in session
//	@Tags			auth
//	@Accept			json
//
//	@Param			request	body		SignInDTO	true	"Sign In DTO"
//
//	@Success		200		{object}	Session
//	@Failure		400		{object}	wool.Error
//	@Failure		422		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/auth/sign-in [post]
func nopAuthSignIn() {}

//	@Summary		OTP
//	@Description	get new session
//	@Tags			auth
//	@Accept			json
//
//	@Param			request	body		OtpDTO	true	"Otp DTO"
//
//	@Success		200		{object}	Session
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		422		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/auth/otp [post]
//	@Security		SysAuth
func nopAuthOTP() {}

//	@Summary		Refresh Token
//	@Description	get new session
//	@Tags			auth
//	@Accept			json
//
//	@Param			request	body		RefreshTokenDTO	true	"Refresh Token DTO"
//
//	@Success		200		{object}	Session
//	@Failure		400		{object}	wool.Error
//	@Failure		422		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/auth/refresh [post]
func nopAuthRefresh() {}

//	@Summary		List feeds
//	@Description	get feeds
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int			false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int			false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Feed	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		403		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/feeds [get]
//	@Security		SysAuth
func nopFeedList() {}

//	@Summary		Show a feed
//	@Description	get feed by ID
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Feed ID"	Format(uuid)
//	@Success		200	{object}	entity.Feed	"OK"
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/feeds/{id} [get]
//	@Security		SysAuth
func nopFeedByID() {}

//	@Summary		Create feed
//	@Description	add new feed
//	@Tags			feeds
//	@Accept			json
//
//	@Param			request	body		CreateFeedDTO	true	"Create Feed DTO"
//
//	@Header			201		{string}	Location		"/feeds/{id}"
//	@Success		201
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/feeds [post]
//	@Security		SysAuth
func nopCreateFeed() {}

//	@Summary		Update feed
//	@Description	edit feed
//	@Tags			feeds
//	@Accept			json
//
//	@Param			id		path	string			true	"Feed ID"	Format(uuid)
//	@Param			request	body	UpdateFeedDTO	true	"Update Feed DTO"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/feeds/{id} [patch]
//	@Security		SysAuth
func nopUpdateFeed() {}

//	@Summary		Delete feed
//	@Description	delete feed
//	@Tags			feeds
//	@Accept			json
//
//	@Param			id	path	string	true	"Feed ID"	Format(uuid)
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/feeds/{id} [delete]
//	@Security		SysAuth
func nopDeleteFeed() {}

//	@Summary		List chats
//	@Description	get chats
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int			false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int			false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Chat	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		403		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/chats [get]
//	@Security		SysAuth
func nopChatList() {}

//	@Summary		Show a chat
//	@Description	get chat by ID
//	@Tags			chats
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Chat ID"	Format(uuid)
//	@Success		200	{object}	entity.Chat	"OK"
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/chats/{id} [get]
//	@Security		SysAuth
func nopChatByID() {}

//	@Summary		Create chat
//	@Description	add new chat
//	@Tags			chats
//	@Accept			json
//
//	@Param			request	body		CreateChatDTO	true	"Create Chat DTO"
//
//	@Header			201		{string}	Location		"/chats/{id}"
//	@Success		201
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/chats [post]
//	@Security		SysAuth
func nopCreateChat() {}

//	@Summary		Update chat
//	@Description	edit chat
//	@Tags			chats
//	@Accept			json
//
//	@Param			id		path	string			true	"Chat ID"	Format(uuid)
//	@Param			request	body	UpdateChatDTO	true	"Update Chat DTO"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/chats/{id} [patch]
//	@Security		SysAuth
func nopUpdateChat() {}

//	@Summary		Delete chat
//	@Description	delete chat
//	@Tags			chats
//	@Accept			json
//
//	@Param			id	path	string	true	"Chat ID"	Format(uuid)
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/chats/{id} [delete]
//	@Security		SysAuth
func nopDeleteChat() {}

//	@Summary		List jobs
//	@Description	get jobs
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int			false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int			false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Job	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		403		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/jobs [get]
//	@Security		SysAuth
func nopJobList() {}

//	@Summary		Show a job
//	@Description	get job by ID
//	@Tags			jobs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Job ID"	Format(uuid)
//	@Success		200	{object}	entity.Job	"OK"
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/jobs/{id} [get]
//	@Security		SysAuth
func nopJobByID() {}

//	@Summary		Create job
//	@Description	add new job
//	@Tags			jobs
//	@Accept			json
//
//	@Param			request	body		CreateJobDTO	true	"Create Job DTO"
//
//	@Header			201		{string}	Location		"/jobs/{id}"
//	@Success		201
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/jobs [post]
//	@Security		SysAuth
func nopCreateJob() {}

//	@Summary		Update job
//	@Description	edit job
//	@Tags			jobs
//	@Accept			json
//
//	@Param			id		path	string			true	"Job ID"	Format(uuid)
//	@Param			request	body	UpdateJobDTO	true	"Update Job DTO"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/jobs/{id} [patch]
//	@Security		SysAuth
func nopUpdateJob() {}

//	@Summary		Delete job
//	@Description	delete job
//	@Tags			jobs
//	@Accept			json
//
//	@Param			id	path	string	true	"Job ID"	Format(uuid)
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/jobs/{id} [delete]
//	@Security		SysAuth
func nopDeleteJob() {}

//	@Summary		List articles
//	@Description	get articles
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int				false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int				false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Article	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		403		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/articles [get]
//	@Security		SysAuth
func nopArticleList() {}

//	@Summary		Show an article
//	@Description	get article by ID
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Article ID"	Format(uuid)
//	@Success		200	{object}	entity.Article	"OK"
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/articles/{id} [get]
//	@Security		SysAuth
func nopArticleByID() {}

//	@Summary		Update article
//	@Description	edit article
//	@Tags			articles
//	@Accept			json
//
//	@Param			id		path	string				true	"Article ID"	Format(uuid)
//	@Param			request	body	UpdateArticleDTO	true	"Update Article DTO"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/articles/{id} [patch]
//	@Security		SysAuth
func nopUpdateArticle() {}

//	@Summary		Delete article
//	@Description	delete article
//	@Tags			articles
//	@Accept			json
//
//	@Param			id	path	string	true	"Article ID"	Format(uuid)
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/articles/{id} [delete]
//	@Security		SysAuth
func nopDeleteArticle() {}
