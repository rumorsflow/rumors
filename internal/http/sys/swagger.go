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

//	@Summary		List sites
//	@Description	get sites
//	@Tags			sites
//	@Accept			json
//	@Produce		json
//	@Param			index	query		int			false	"Page Index"	default(0)	minimum(0)
//	@Param			size	query		int			false	"Page Size"		default(20)	minimum(1)	maximum(100)
//	@Success		200		{array}		entity.Site	"OK"
//	@Failure		400		{object}	wool.Error
//	@Failure		401		{object}	wool.Error
//	@Failure		403		{object}	wool.Error
//	@Failure		500		{object}	wool.Error
//	@Router			/sites [get]
//	@Security		SysAuth
func nopSiteList() {}

//	@Summary		Show a site
//	@Description	get site by ID
//	@Tags			sites
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string		true	"Site ID"	Format(uuid)
//	@Success		200	{object}	entity.Site	"OK"
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/sites/{id} [get]
//	@Security		SysAuth
func nopSiteByID() {}

//	@Summary		Create site
//	@Description	add new site
//	@Tags			sites
//	@Accept			json
//
//	@Param			request	body		CreateSiteDTO	true	"Create Site DTO"
//
//	@Header			201		{string}	Location		"/sites/{id}"
//	@Success		201
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/sites [post]
//	@Security		SysAuth
func nopCreateSite() {}

//	@Summary		Update site
//	@Description	edit site
//	@Tags			sites
//	@Accept			json
//
//	@Param			id		path	string			true	"Site ID"	Format(uuid)
//	@Param			request	body	UpdateSiteDTO	true	"Update Site DTO"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		409	{object}	wool.Error
//	@Failure		422	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/sites/{id} [patch]
//	@Security		SysAuth
func nopUpdateSite() {}

//	@Summary		Delete site
//	@Description	delete site
//	@Tags			sites
//	@Accept			json
//
//	@Param			id	path	string	true	"Site ID"	Format(uuid)
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/sites/{id} [delete]
//	@Security		SysAuth
func nopDeleteSite() {}

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

//	@Summary		Delete queue
//	@Description	delete queue
//	@Tags			queues
//
//	@Param			qname	path	string	true	"Queue name"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		404	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/queues/{qname} [delete]
//	@Security		SysAuth
func nopDeleteQueue() {}

//	@Summary		Pause queue
//	@Description	pause queue
//	@Tags			queues
//
//	@Param			qname	path	string	true	"Queue name"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/queues/{qname}/pause [post]
//	@Security		SysAuth
func nopPauseQueue() {}

//	@Summary		Resume queue
//	@Description	resume queue
//	@Tags			queues
//
//	@Param			qname	path	string	true	"Queue name"
//
//	@Success		204
//	@Failure		400	{object}	wool.Error
//	@Failure		401	{object}	wool.Error
//	@Failure		403	{object}	wool.Error
//	@Failure		500	{object}	wool.Error
//	@Router			/queues/{qname}/resume [post]
//	@Security		SysAuth
func nopResumeQueue() {}
