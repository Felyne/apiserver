import sys
import requests
from json import dumps


BASE_URL = 'http://127.0.0.1:8080'

tokenHeader = None


def getUrl(path):
    return BASE_URL + path


def login():
    global tokenHeader
    url = getUrl('/login')
    req = {
        'username': 'admin',
        'password': 'admin'
    }
    result = requests.post(url, json=req).json()
    if result['code'] != 0:
        print(result)
        sys.exit(1)
    tokenHeader = {'Authorization': 'Bearer ' + result['data']['token']}


def createUser():
    url = getUrl('/v1/user')
    req = {
        'username': 'test',
        'password': '123456'
    }
    result = requests.post(url, json=req, headers=tokenHeader).json()
    print(dumps(result, indent=2))


def getUser():
    url = getUrl('/v1/user/test')
    userList = requests.get(url, headers=tokenHeader).json()
    print(dumps(userList, indent=2))


def listUser():
    url = getUrl('/v1/user')
    userList = requests.get(url, headers=tokenHeader).json()
    print(dumps(userList, indent=2))


def main():
    login()
    createUser()
    getUser()
    listUser()


if __name__ == '__main__':
    main()
