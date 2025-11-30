import { createRouter, createWebHistory } from 'vue-router'
import NProgress from 'nprogress'
import Layout from '../layout/Layout.vue'
import AdminLayout from '../layout/AdminLayout.vue'
import { useSiteConfig, initSiteConfig } from '../store/site'

const routes = [
  // 前台路由
  {
    path: '/',
    component: Layout,
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'article/:id',
        name: 'ArticleDetail',
        component: () => import('../views/ArticleDetail.vue'),
        meta: { title: '文章详情' }
      },
      {
        path: 'category',
        name: 'Category',
        component: () => import('../views/Category.vue'),
        meta: { title: '分类' }
      },
      {
        path: 'category/:slug',
        name: 'CategoryDetail',
        component: () => import('../views/CategoryDetail.vue'),
        meta: { title: '分类详情' }
      },
      {
        path: 'tag',
        name: 'Tag',
        component: () => import('../views/Tag.vue'),
        meta: { title: '标签' }
      },
      {
        path: 'tag/:slug',
        name: 'TagDetail',
        component: () => import('../views/TagDetail.vue'),
        meta: { title: '标签详情' }
      },
      {
        path: 'topic',
        name: 'Topic',
        component: () => import('../views/Topic.vue'),
        meta: { title: '专题' }
      },
      {
        path: 'topic/:id',
        name: 'TopicDetail',
        component: () => import('../views/TopicDetail.vue'),
        meta: { title: '专题详情' }
      },
      {
        path: 'about',
        name: 'About',
        component: () => import('../views/About.vue'),
        meta: { title: '关于' }
      },
      {
        path: 'my-articles',
        name: 'MyArticles',
        component: () => import('../views/MyArticles.vue'),
        meta: { title: '我的文章', requiresAuth: true }
      },
      {
        path: 'author/:id',
        name: 'AuthorArticles',
        component: () => import('../views/AuthorArticles.vue'),
        meta: { title: '作者文章' }
      }
    ]
  },
  // 登录
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  // 新版登录页面
  {
    path: '/login-new',
    name: 'LoginNew',
    component: () => import('../views/LoginNew.vue'),
    meta: { title: '登录（新版）' }
  },
  // 注册页面
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { title: '注册' }
  },
  // OAuth 回调页面（处理OAuth登录成功后的回调）
  {
    path: '/oauth/:provider/callback',
    name: 'OAuthCallback',
    component: () => import('../views/OAuthCallback.vue'),
    meta: { title: '登录回调' }
  },
  // 后台路由
  {
    path: '/admin',
    component: AdminLayout,
    redirect: '/admin/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/admin/Dashboard.vue'),
        meta: { title: '仪表盘' }
      },
      {
        path: 'articles',
        name: 'ArticleManage',
        component: () => import('../views/admin/ArticleManage.vue'),
        meta: { title: '文章管理' }
      },
      {
        path: 'articles/edit/:id?',
        name: 'ArticleEdit',
        component: () => import('../views/admin/ArticleEdit.vue'),
        meta: { title: '编辑文章' }
      },
      {
        path: 'category',
        name: 'CategoryManage',
        component: () => import('../views/admin/CategoryManage.vue'),
        meta: { title: '分类管理' }
      },
      {
        path: 'tags',
        name: 'TagManage',
        component: () => import('../views/admin/TagManage.vue'),
        meta: { title: '标签管理' }
      },
      {
        path: 'topics',
        name: 'TopicManage',
        component: () => import('../views/admin/TopicManage.vue'),
        meta: { title: '专题管理' }
      },
      {
        path: 'users',
        name: 'UserManage',
        component: () => import('../views/admin/UserManage.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'roles',
        name: 'RoleManage',
        component: () => import('../views/admin/RoleManage.vue'),
        meta: { title: '角色管理' }
      },
      {
        path: 'permissions',
        name: 'PermissionManage',
        component: () => import('../views/admin/PermissionManage.vue'),
        meta: { title: '权限管理' }
      },
      {
        path: 'menus',
        name: 'MenuManage',
        component: () => import('../views/admin/MenuManage.vue'),
        meta: { title: '菜单管理' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/admin/Settings.vue'),
        meta: { title: '系统设置' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/admin/Profile.vue'),
        meta: { title: '个人中心' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    // 如果有保存的位置（浏览器前进/后退），则返回到该位置
    if (savedPosition) {
      return savedPosition
    } else {
      // 否则滚动到页面顶部
      return { top: 0, behavior: 'smooth' }
    }
  }
})

// NProgress Configuration
NProgress.configure({ showSpinner: false })

router.beforeEach((to, from, next) => {
  NProgress.start()
  const { siteName } = useSiteConfig()
  const pageTitle = to.meta.title || '首页'
  const siteTitle = siteName.value || 'Blog System'
  document.title = `${pageTitle} - ${siteTitle}`
  
  // 检查是否需要登录
  const token = localStorage.getItem('token')
  const isAdminRoute = to.path.startsWith('/admin')
  const isLoginRoute = to.path === '/login' || to.path === '/login-new' || to.path === '/register' || to.path.startsWith('/oauth/')
  const requiresAuth = to.meta.requiresAuth
  
  // 如果需要登录但未登录，重定向到登录页
  if (requiresAuth && !token) {
    next({
      path: '/login-new',
      query: { redirect: to.fullPath } // 保存原始路径，登录后可以跳转回来
    })
    return
  }
  
  // 如果是管理后台路由且未登录，重定向到登录页
  if (isAdminRoute && !token) {
    next({
      path: '/login-new',
      query: { redirect: to.fullPath } // 保存原始路径，登录后可以跳转回来
    })
    return
  }
  
  // 如果已登录且访问登录页，重定向到管理后台首页
  if (isLoginRoute && token) {
    next('/admin/dashboard')
    return
  }
  
  next()
})

router.afterEach(() => {
  NProgress.done()
})

export default router
