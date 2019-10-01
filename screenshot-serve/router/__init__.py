#!/usr/bin/env python3
# -*- coding: utf-8 -*-
import os
import time
import json
import requests
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.by import By
from flask import Blueprint, request

router = Blueprint('router', __name__)

@router.route('/')
def index():
    return 'hello world'

# 截图
@router.route('/screenshot', methods=['POST'])
def screenshot():
  data = request.get_data()
  json_data = json.loads(data.decode("utf-8"))
  previewUrl = json_data.get('previewUrl')
  callbackUrl = json_data.get('callbackUrl')
  
  print("json_data", json_data)
  print("previewUrl", previewUrl)
  if previewUrl is None:
      return

  opt = Options()
  opt.add_argument('--headless') # 不显示浏览器

  driver = webdriver.Chrome(options=opt)              # 打开Google浏览器
  driver.get(previewUrl)

  locator = (By.ID, "render-complete")

  # 智能等待10s之后获取元素
  def find_elements(driver, locator):
      WebDriverWait(driver, 10).until(EC.presence_of_element_located(locator))

  try:
    print(os.getcwd())
    file_path = os.getcwd() + "/static/images"
    if not os.path.exists(file_path):
        os.makedirs(file_path)
        print("目录新建成功：%s" % file_path)
  except BaseException as msg:
      print("新建目录失败：%s" % msg)

  try:
      find_elements(driver, locator)
      ticks = time.time()
      file_name = '/%s.png'%(ticks)
      driver.save_screenshot(file_path + file_name)
      print('save_screenshot')
  except:
      print("ele can't find")
  finally:
      driver.quit()

  data = json.dumps({
      'token': json_data.get('token'),
      'elementId': json_data.get('elementId'),
      'projectId': json_data.get('projectId'),
      'url' : 'http://localhost:5000/static/images' + file_name
  })    

  if callbackUrl is not None:
      r = requests.post(callbackUrl, data)
      print("response")
      print(r.content)

  return data