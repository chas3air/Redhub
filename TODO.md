Names of branches:
    1. task-usersManageService: creating UsersManageService
    2. task-auth: creating Auth
    3. task-api-gateway: creating Api-Gateway
    4. fix-proto: fix-proto in a-g, auth, ums
    5. task-add-usersdb
    6. task-add-description: add description to user struct
    7. task-articlesManageService: creaitng ArticlesManageService
    8. fix-secret-store: fix storing clients secrets
    9. task-addArticleClient: add article client to article service
    10. fix-UserManageService: fix users client and usersmanageservice
    11. task-commentManageService: create CommentsManageService
    12. task-addCommentToAG: add commentManService to api-gateway
    13. fix-UAprotos: fix users and articles proto to return obj in insert and update
    14. fix-auth: fix auth, adding access and refresh tokens, maybe add cache
    15. task-web: do web
    15. task-stats: add statistics to app


Tasks:
    нужно сделать базу данных которая будет хранить клиентов их id и секреты. сделать какую-нибудь простую апиху которая регает их. таску оставлю на завтра. Тогда же напишу мидлваер на валидацию зареганый пользователей и нет, и раскидаю права.