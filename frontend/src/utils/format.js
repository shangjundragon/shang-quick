// 格式化毫秒时间戳为中文本地时间字符串（yyyy/MM/dd HH:mm:ss）
export function formatTimestamp(timestamp) {
  if (!timestamp) return ''
  const date = new Date(Number(timestamp))
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
