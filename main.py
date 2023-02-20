from os import environ
from check import CheckIn
from push import send_msg_serverJ, send_msg_pushplus, send_bark
from config import config


def main():
    conf = config()
    # 获取actions secrets配置的cookie SendKey
    ck = conf[0]['cookies']
    # send_key = environ.get('SendKey')
    # token = environ.get('token')
    bark_url = conf[0]['bark_url']
    try:
        title, text = CheckIn(ck)
        print('签到成功！')
    except Exception as e:
        print('程序出错！')
        title = '程序出错！'
        text = e
    finally:
        # print(title)
        print(text)
        # Text = Text.replace('\n', '%0D%0A%0D%0A')

        # if send_key != '':
        #     rsp = send_msg_serverJ(send_key, title, text)  # 推送消息，无SendKey不推送
        #     print(rsp)
        #
        # if token != '':
        #     rsp = send_msg_pushplus(token, title, text)  # 推送消息，无token不推送
        #     print(rsp)

        if bark_url != '':
            rsp = send_bark(bark_url, title, text)
            print(rsp)


if __name__ == '__main__':
    main()
