# coding:utf-8
import yaml


def config():
    # 获取yaml文件路径
    yaml_path = 'config.yml'

    with open(yaml_path, 'rb') as f:
        # yaml文件通过---分节，多个节组合成一个列表
        date = yaml.safe_load_all(f)
        # print(list(date))
        return list(date)
