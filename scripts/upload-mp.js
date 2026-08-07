// 命令行上传微信小程序体验版 / 生成预览二维码，脱离微信开发者工具
// 用法：node scripts/upload-mp.js --desc "版本说明" [--preview] [--key mp/upload-key.pem] [--robot 1]
const path = require('path')
const ci = require('miniprogram-ci')
const manifest = require('../mp/src/manifest.json')

function parseArgs(argv) {
  const args = { preview: false, key: 'mp/upload-key.pem', robot: 1, desc: '' }
  for (let i = 0; i < argv.length; i++) {
    const arg = argv[i]
    if (arg === '--preview') args.preview = true
    else if (arg === '--desc') args.desc = argv[++i]
    else if (arg === '--key') args.key = argv[++i]
    else if (arg === '--robot') args.robot = Number(argv[++i])
  }
  return args
}

async function main() {
  const args = parseArgs(process.argv.slice(2))
  if (!args.desc) {
    console.error('缺少 --desc 参数，请提供版本说明，例如：node scripts/upload-mp.js --desc "修复播放页样式"')
    process.exit(1)
  }

  const project = new ci.Project({
    appid: manifest.appid,
    type: 'miniProgram',
    projectPath: path.resolve(__dirname, '../mp/dist/build/mp-weixin'),
    privateKeyPath: path.resolve(__dirname, '..', args.key),
    ignores: ['node_modules/**/*'],
  })

  const setting = { es6: true, minify: true }

  if (args.preview) {
    await ci.preview({
      project,
      version: manifest.versionName,
      desc: args.desc,
      setting,
      robot: args.robot,
      qrcodeFormat: 'terminal',
    })
    console.log('✅ 预览二维码已生成，请扫码体验')
  } else {
    const uploadResult = await ci.upload({
      project,
      version: manifest.versionName,
      desc: args.desc,
      setting,
      robot: args.robot,
    })
    console.log('✅ 上传成功，已生成体验版本：', uploadResult)
  }
}

main().catch((err) => {
  console.error('❌ 操作失败：', err.message || err)
  process.exit(1)
})
