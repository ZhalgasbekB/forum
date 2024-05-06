## **forum** 

### **Objectives**

The project implements the functionality of a web resource in the form of a forum. The forum has subforums dedicated to transport topics. At this stage of the project, a sub-forum has been developed dedicated to suburban transportation along the section of the Moscow Central Diameter No. 3 (named D3).

Functionally, the forum has the ability to register new users with further login to use the web resource. The forum home page contains information about current topics being discussed, comments and reactions. Each user can leave a post with a topic in one of the following categories:

- Trains
- Stations
- Schedule
- Rates
- Other

Each user has the opportunity to put reactions to posts and comments, as well as make changes and delete them.

All the user's liked posts are found when you go to "Like Topics".

Themes that were created by the user are available when you go to "My Themes".

Upon completion, the user can end the session by clicking on “Logout”.

The project is built on a microservice architecture, which includes two services:
- Client
- Server

The project uses the following technologies:
- Go programming language
- hypertext markup language HTML
- CSS styling language
- for building Docker service components

Interaction is based on REST API principles.

### **Instructions**

Procedure for the user:
<br>

1. Clone the project from the repository
2. Go to the project root folder
3. Enter the ` docker-compose up ` command
4. In the command line of the client terminal, wait for the message about the successful launch of the service "The system check was successful, the client server has started"
5. Go to [http://localhost:8082/](http://localhost:8082/)


### **Tests**

To test, go to the root folder of the project and run the command: ` go test ./...`

### **Autors**

[@zhbolatov](https://01.alem.school/git/zhbolatov)
[@lzhuk](https://01.alem.school/git/lzhuk)
[@dbaitako](https://01.alem.school/git/dbaitako)
