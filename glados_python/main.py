#!/usr/bin/env python
# coding=utf-8
from pytz import timezone
from json import dumps
from requests import post, get
from datetime import datetime
import yaml


def send_bark(url, title, text):
    if not url:
        return 'bark: 未配置，无法进行消息推送.'
    print('=================================================================\nBark: 开始推送消息！')
    uri = url + '/' + title + '/' + text
    rsp = get(uri)
    return rsp.json()['message']


def checkin(cookie):
    url = "https://glados.rocks/api/user/checkin"
    url2 = "https://glados.rocks/api/user/status"
    referer = 'https://glados.rocks/console/checkin'
    origin = "https://glados.rocks"
    useragent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"
    payload = {"token":"glados.one"}
    check_in = post(
        url,
        headers={
            'cookie': cookie,
            'referer': referer,
            'origin': origin,
            'user-agent': useragent,
            'content-type': 'application/json;charset=UTF-8'
        },
        data=dumps(payload)
    )
    state = get(
        url2,
        headers={
            'cookie': cookie,
            'referer': referer,
            'origin': origin,
            'user-agent': useragent
        }
    )
    tz = timezone('Asia/Shanghai')
    time_now = str(datetime.now(tz=tz))[:19]

    # 重复打卡 code==1
    if check_in.json()['code'] == -2:
        raise Exception(check_in.json()['message'])

    mess = check_in.json()['message']
    time = state.json()['data']['leftDays']
    email = state.json()['data']['email']
    days = time.split('.')[0]
    msg = f'现在时间是：{time_now}\nemail: {email}\ncheckin: {check_in.status_code} | state: {state.status_code}\n{mess}\n剩余天数：{days}天'

    check_in.close()
    state.close()

    return f'{mess}，剩余{days}天', msg


def main():
    # 获取yaml文件路径
    yaml_path = 'config.yml'
    # 使用open()函数读取config.yaml文件
    yaml_file = open(yaml_path, "r", encoding="utf-8")
    # 读取文件中的内容
    file_data = yaml_file.read()
    yaml_file.close()

    # 加载数据流，返回字典类型数据
    conf = yaml.load(file_data, Loader=yaml.FullLoader)

    ck = conf['cookies']
    bark_url = conf['bark_url']

    for c in ck:
        try:
            title, text = checkin(c)
            print('签到成功！')
        except Exception as e:
            print('程序出错！')
            title = '程序出错！'
            if len(e.args) > 0:
                text = '网络信息: ' + e.args[0]
            else:
                text = "没有获取到网络信息"
        finally:
            # print(title)
            print(text)
            # Text = Text.replace('\n', '%0D%0A%0D%0A')

            if bark_url != '':
                rsp = send_bark(bark_url, title, text)
                print(rsp)


if __name__ == '__main__':
    main()