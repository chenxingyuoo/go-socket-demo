<template>
  <div id="app">
    <Stage></Stage>
  </div>
</template>

<script>
import Stage from './components/stage/index'

export default {
  name: 'app',
  components: {
    Stage
  }
}
</script>

<style lang="less">
*{
  margin: 0;
  box-sizing: border-box;
}
html, body, #app {
  width: 100%;
  height: 100%;
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
