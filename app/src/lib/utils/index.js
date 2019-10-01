/**
 * 获取2个平面的宽比、高比
 * @param  {object}  originRect  原宽高对象
 * @param  {object}  targetRect  目标宽高对象
 * @return {object}              宽高比
 * @example
 * getRect2WHRatio({width: 100, height: 100}, {width: 200, height: 200})
 * // => {heightRatio: 2,maxRatio: 2, minRatio: 2, widthRatio: 2}
 */
export const getRect2WHRatio = (originRect, targetRect) => {
  const widthRatio = targetRect.width / originRect.width
  const heightRatio = targetRect.height / originRect.height
  return {
    widthRatio: widthRatio,
    heightRatio: heightRatio,
    maxRatio: Math.max(widthRatio, heightRatio),
    minRatio: Math.min(widthRatio, heightRatio)
  }
}

/**
 * 获取唯一id
 */
export const guid = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0; var v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}
