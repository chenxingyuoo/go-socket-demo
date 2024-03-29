<template>
  <div id="stage-box">
    <div id="render-complete"></div>
    <!-- element -->
    <stage-element
      v-for="(item, index) in data"
      ref="element-item"
      :data="item"
      :index="index"
      :key="index"
      @onMousedown="handleMousedown"
      ></stage-element>

      <!-- resize -->
      <vue-drag-rotate-resize
        v-show="Object.keys(resizeObject).length"
        :top="resizeObject.y || 0"
        :left="resizeObject.x || 0"
        :width="resizeObject.width || 0"
        :height="resizeObject.height || 0"
        :rotateAngle="resizeObject.rotate || 0"
        :rotatable="!isMultiple"
        :aspectRatio="isMultiple ? resizeObject.width / resizeObject.height : 0"
        :onDragStart="handleDragStart"
        :onDrag="handleDrag"
        :onDragEnd="handleDragEnd"
        :onResizeStart="handleResizeStart"
        :onResize="handleResize"
        :onResizeEnd="hanlderResizeEnd"
        :onRotateStart="handleRotateStart"
        :onRotate="handleRotate"
        :onRotateEnd="handleRotateEnd"
      >
      <template slot="rotate">
        <div class="rotate-btn"></div>
      </template>
      </vue-drag-rotate-resize>

      <button class="save-btn" @click="saveElement">保存</button>

      <div class="thumbnail" v-if="thumbnailSrc">
         <img :src="thumbnailSrc">
      </div>
  </div>
</template>

<script>
import { guid, getRect2WHRatio } from '@/lib/utils/'
import cookie from '@/lib/utils/cookie'
import cloneDeep from 'lodash/cloneDeep'
import StageElement from '../stage-element/index'
import bg from '../../assets/bg.jpg'
import btn from '../../assets/btn.png'

if (!cookie.get('token')) {
  cookie.set('token', guid())
}

const token = cookie.get('token')

export default {
  data () {
    this.websocket = null
    return {
      isDraging: false,
      resizeObject: {},
      activeElement: null,
      activedElements: [],
      thumbnailSrc: '',
      data: [
        {
          x: 100,
          y: 60,
          width: 200,
          height: 100,
          rotate: 0,
          zIndex: 100,
          src: bg
        },
        {
          x: 120,
          y: 180,
          width: 162.5,
          height: 39,
          rotate: 0,
          zIndex: 101,
          src: btn
        }
      ]
    }
  },
  components: {
    StageElement
  },
  computed: {
    isMultiple () {
      return this.activedElements.length > 1
    }
  },
  created () {
    this.wsInit()

    document.addEventListener('click', (e) => {
      if (this.getDraging()) {
        return
      }

      if (!e.target.closest('#stage-box')) {
        this.hideEdit()
      }
    }, false)
  },
  methods: {
    wsInit () {
      const wsUri = `ws://localhost:8000/ws/${token}`
      const websocket = new WebSocket(wsUri)
      websocket.onopen = (evt) => {
        console.log('Connection open ...')

        websocket.send(JSON.stringify({
          type: 'testMessage',
          data: 'test data'
        }))
      }
      websocket.onclose = (evt) => {
        console.log('Connection close !!!')
        alert('Connection close !!!')
      }
      websocket.onmessage = (evt) => {
        console.log('onmessage', evt)
        try {
          const jsonData = JSON.parse(evt.data)
          if (jsonData.type === 'thumbnail') {
            this.thumbnailSrc = jsonData.data.url
          }
        } catch (e) {
          console.log(e)
        }
      }
      websocket.onerror = (evt) => {
        console.log('Connection onerror !!!')
      }

      this.websocket = websocket
    },
    hideEdit () {
      this.resizeObject = {}
    },
    setDraging () {
      this.isDraging = true
    },
    cancelDraging () {
      setTimeout(() => {
        this.isDraging = false
      }, 10)
    },
    getDraging () {
      return this.isDraging
    },
    getRect () {
      let object = {
        x: [],
        y: [],
        width: [],
        height: []
      }

      // get x、y
      this.activedElements.map((element, index) => {
        object.x.push(element.x)
        object.y.push(element.y)
      })
      const minX = Math.min(...object.x)
      const minY = Math.min(...object.y)

      // get width、height
      this.activedElements.map((element, index) => {
        const { right, bottom } = this.$refs['element-item'][index]['$el'].getBoundingClientRect()
        object.width.push(right - minX)
        object.height.push(bottom - minY)
      })

      return {
        x: minX,
        y: minY,
        width: Math.max(...object.width),
        height: Math.max(...object.height)
      }
    },
    handleMousedown (e, data, index) {
      if (e.metaKey || e.ctrlKey) {
        this.activedElements.push(data)
        if (this.activedElements.length === 1) {
          this.resizeObject = cloneDeep(data)
        } else {
          const rect = this.getRect()
          this.resizeObject = rect
        }
      } else {
        this.resizeObject = cloneDeep(data)
        this.activedElements = [data]
      }
    },
    handleDragStart () {
      this.setDraging()
    },
    handleDrag (e, deltaX, deltaY) {
      this.resizeObject.x += deltaX
      this.resizeObject.y += deltaY

      this.activedElements.map((element, index) => {
        element.x += deltaX
        element.y += deltaY
      })
    },
    handleDragEnd () {
      this.cancelDraging()
    },
    handleResizeStart () {
      this.setDraging()
      this.oldResizeObject = cloneDeep(this.resizeObject)
      this.oldActivedElements = cloneDeep(this.activedElements)
    },
    handleResize (e, { top, left, width, height }, isShiftKey, type) {
      const options = {
        y: Math.round(top),
        x: Math.round(left),
        width: Math.round(width),
        height: Math.round(height)
      }

      this.resizeRatioObject = getRect2WHRatio(options, this.oldResizeObject)

      const { widthRatio, heightRatio } = this.resizeRatioObject

      const diffX = options.x - this.oldResizeObject.x
      const diffY = options.y - this.oldResizeObject.y

      const rect = this.getRect()

      this.activedElements.map((element, index) => {
        const oldElement = this.oldActivedElements[index]

        const width = Math.round(oldElement.width / widthRatio)
        const height = Math.round(oldElement.height / heightRatio)

        const xRatio = (oldElement.x - this.oldResizeObject.x) / this.oldResizeObject.width
        const yRatio = (oldElement.y - this.oldResizeObject.y) / this.oldResizeObject.height

        const x = diffX + this.oldResizeObject.x + Math.round(options.width * xRatio)
        const y = diffY + this.oldResizeObject.y + Math.round(rect.height * yRatio)

        Object.assign(element, {
          x,
          y,
          width,
          height
        })
      })

      Object.assign(this.resizeObject, options)
    },
    hanlderResizeEnd () {
      this.cancelDraging()
      if (this.isMultiple) {
        const rect = this.getRect()
        Object.assign(this.resizeObject, rect)
      }
    },
    handleRotateStart () {
      this.setDraging()
    },
    handleRotate (e, rotateAngle) {
      this.resizeObject.rotate = rotateAngle
      this.activedElements.map((element, index) => {
        element.rotate = rotateAngle
      })
    },
    handleRotateEnd () {
      this.cancelDraging()
    },
    saveElement () {
      for (let index = 0; index < 1; index++) {
        fetch('/api/element/save', {
          method: 'POST',
          body: JSON.stringify({
            projectId: 1,
            elementId: 38
          }),
          headers: {
            token: token
          }
        }).then((res) => {
          return res.json()
        }).then((json) => {
          console.log('json', json)
          alert('保存成功')
        }).catch(() => {

        })
      }
    }
  }
}
</script>

<style lang="less" scoped>

#stage-box{
  position: relative;
  width: 100%;
  height: 100%;
}
.rotate-btn {
  position: absolute;
  width: 14px;
  height: 14px;
  left: 50%;
  top: -32px;
  margin-left: -7px;
  border: 1px solid #409EFF;
  border-radius: 50%;
  box-sizing: border-box;
}
.save-btn{
  position: absolute;
  top: 0;
  right: 0;
}

.thumbnail {
  width: 120px;
  height: 120px;
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  bottom: 0;
  border: 1px solid #ddd;
  img {
    max-width: 100%;
  }
}
</style>
