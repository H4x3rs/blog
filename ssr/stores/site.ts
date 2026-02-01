import { defineStore } from 'pinia'

interface SiteState {
  siteName: string
  bannerTitle: string
  bannerSubtitle: string
  icpNumber: string
  loaded: boolean
}

export const useSiteStore = defineStore('site', {
  state: (): SiteState => ({
    siteName: 'Blog System',
    bannerTitle: '分享编程心得\n记录技术成长',
    bannerSubtitle: '探索技术之美，记录成长足迹。分享 Go, Vue, 云原生等前沿技术心得。',
    icpNumber: '',
    loaded: false
  }),

  actions: {
    async loadConfig() {
      if (this.loaded) return

      try {
        const { getBanner } = useSettingsApi()
        const data = await getBanner()
        
        if (data) {
          this.siteName = data.siteName || this.siteName
          this.bannerTitle = data.bannerTitle || this.bannerTitle
          this.bannerSubtitle = data.bannerSubtitle || this.bannerSubtitle
          this.icpNumber = data.icpNumber || ''
        }
        
        this.loaded = true
      } catch (error) {
        console.error('加载网站配置失败:', error)
      }
    }
  }
})
