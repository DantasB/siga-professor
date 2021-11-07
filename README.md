# Siga-Professor

![demonstration](https://cdn.discordapp.com/attachments/539836343094870016/906955024108752988/unknown.png)

## Table of Contents

<!--ts-->
   * [About](#about)
   * [Requirements](#requirements)
   * [How to use](#how-to-use)
      * [API Routes](#api-routes)
      * [Setting up Program](#program-setup)
   * [Technologies](#technologies)
<!--te-->

## About

This repository is a Simple api that access every UFRJ course and get all the oppened disciplines and their professors. It was built using golang.

## Requirements

To run this repository by yourself you will need to have the go compiler on your operational system to run the code bellow.

## How to use

### API Routes

- *'/filldatabase':* This route access every UFRJ course and stores the data in a mongodb.
- *'/professors/<professor_name>':* This route only works after the first one.

### Program Setup

```bash
# Clone this repository
$ git clone <https://github.com/DantasB/Siga-Professor>

# Access the project page on your terminal
$ cd Siga-Professor/

# Create a .env file
$ touch .env  

# Create the following parameters
 CONNECTION_URL="mongodb+srv://{USER}:{PASSWORD}@{HOST}/{DATABASE}?{PARAMS}" #Your MongoDB connection url
 MONGODB_USER #Your MongoDB connection username
 MONGODB_PASSWORD #Your MongoDB connection password
 MONGODB_HOST #The database host
 MONGODB_PARAMS #The params of the connection (retryWrites,w, etc...)
 MONGODB_DATABASE #Your database name

# Compile the entire project
$ go run .

# The code will start in the port 8000
```
![demonstration](https://cdn.discordapp.com/attachments/539836343094870016/906954402609397760/unknown.png)


## Technologies

* Golang
* gorilla
