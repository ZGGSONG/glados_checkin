# coding:utf-8
import yaml


# https://blog.csdn.net/weixin_41010198/article/details/111591030
def config():
    # 获取yaml文件路径
    yaml_path = 'config.yml'
    # 使用open()函数读取config.yaml文件
    yaml_file = open(yaml_path, "r", encoding="utf-8")
    # 读取文件中的内容
    file_data = yaml_file.read()
    # print(f"file_date type: {type(file_data)}\nfile_date value:\n{file_data}")
    yaml_file.close()

    # 加载数据流，返回字典类型数据
    y = yaml.load(file_data, Loader=yaml.FullLoader)
    return y
