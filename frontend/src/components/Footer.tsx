import Link from 'next/link';

export default function Footer() {
  return (
    <footer className="border-t border-hairline bg-canvas py-8 pb-20 md:pb-8">
      <div className="max-w-container mx-auto px-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-6 mb-8">
          <div>
            <h4 className="text-sm font-bold mb-3">课程分类</h4>
            <nav className="flex flex-col gap-2">
              <Link href="/courses" className="text-sm text-mute hover:text-primary">全部课堂</Link>
              <Link href="/courses?categoryId=1" className="text-sm text-mute hover:text-primary">AIGC课堂</Link>
              <Link href="/courses?categoryId=2" className="text-sm text-mute hover:text-primary">Blender课堂</Link>
            </nav>
          </div>
          <div>
            <h4 className="text-sm font-bold mb-3">帮助中心</h4>
            <nav className="flex flex-col gap-2">
              <Link href="/vip" className="text-sm text-mute hover:text-primary">成为会员</Link>
              <Link href="/guestbook" className="text-sm text-mute hover:text-primary">留言反馈</Link>
            </nav>
          </div>
          <div>
            <h4 className="text-sm font-bold mb-3">关于</h4>
            <nav className="flex flex-col gap-2">
              <span className="text-sm text-mute">知猿课堂</span>
              <span className="text-sm text-mute">优质网课资源平台</span>
            </nav>
          </div>
          <div>
            <h4 className="text-sm font-bold mb-3">联系方式</h4>
            <nav className="flex flex-col gap-2">
              <span className="text-sm text-mute">客服微信：zioran888</span>
            </nav>
          </div>
        </div>
        <div className="border-t border-hairline pt-4 text-center">
          <span className="text-xs text-mute">© 2024 知猿课堂 All Rights Reserved</span>
        </div>
      </div>
    </footer>
  );
}
