// 播放页互动按钮图标（内联 SVG data URI，替代低保真 emoji，风格统一的线性/面性图标）
function svg(path, fill = '#fff') {
  const body = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48"><path fill="${fill}" d="${path}"/></svg>`
  return 'data:image/svg+xml;utf8,' + encodeURIComponent(body)
}

// 点赞（实心大拇指）
const LIKE_PATH = 'M4 20h6v22H4a2 2 0 0 1-2-2V22a2 2 0 0 1 2-2Zm10.5 0h17.9c1.9 0 3.3 1.8 2.8 3.6l-3.4 13.6a4 4 0 0 1-3.9 3H16a2 2 0 0 1-2-2V21a1 1 0 0 1 .5-1Zm3-16c3 0 4 2.6 4 5.3 0 1.7-.6 3.3-1.4 4.7h6.4c2.9 0 5 2.8 4.2 5.6l-.2.7H14.5V8.3C14.5 5.8 15.6 4 17.5 4Z'
// 分享（转发箭头）
const SHARE_PATH = 'M34 6 46 18 34 30v-7c-11 0-18 3.5-24 15-1-14 7-24 24-25V6Z'
// 收藏（实心心形）
const FAVORITE_PATH = 'M24 42.7l-2.9-2.64C10.8 30.72 4 24.56 4 17c0-6.16 4.84-11 11-11 3.48 0 6.82 1.62 9 4.18C26.18 7.62 29.52 6 33 6c6.16 0 11 4.84 11 11 0 7.56-6.8 13.72-17.1 23.06L24 42.7Z'
// 评论（对话气泡）
const COMMENT_PATH = 'M6 8h36a2 2 0 0 1 2 2v20a2 2 0 0 1-2 2H18l-10 9V32H6a2 2 0 0 1-2-2V10a2 2 0 0 1 2-2Z'
// 倍速（时钟指针风格）
const SPEED_PATH = 'M24 4a20 20 0 1 0 20 20A20 20 0 0 0 24 4Zm2 20h10v4H22V12h4Z'

export const icons = {
  like: svg(LIKE_PATH),
  likeActive: svg(LIKE_PATH, '#409EFF'),
  share: svg(SHARE_PATH),
  favorite: svg(FAVORITE_PATH),
  favoriteActive: svg(FAVORITE_PATH, '#ff4d4f'),
  comment: svg(COMMENT_PATH),
  speed: svg(SPEED_PATH),
}
