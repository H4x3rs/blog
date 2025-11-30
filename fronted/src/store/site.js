import { ref } from 'vue'
import { getBanner } from '@/api/settings'

// 全局网站配置状态
const siteName = ref('Blog System')
const bannerTitle = ref('分享编程心得\n记录技术成长')
const bannerSubtitle = ref('分享编程心得，记录技术成长。分享 Go, Vue, 云原生等前沿技术心得。')

// 初始化网站配置
export const initSiteConfig = async () => {
  try {
    const res = await getBanner()
    if (res.siteName) {
      siteName.value = res.siteName
    }
    if (res.title) {
      bannerTitle.value = res.title
    }
    if (res.subtitle) {
      bannerSubtitle.value = res.subtitle
    }
  } catch (error) {
    console.error('加载网站配置失败:', error)
    // 使用默认值
  }
}

export const useSiteConfig = () => {
  return {
    siteName,
    bannerTitle,
    bannerSubtitle
  }
}




