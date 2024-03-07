import Image from 'next/image';
import Router from 'next/router';
import '@styles/globals.css'
import { HeaderProps } from './Header.type';

const Header = (props: HeaderProps) => {
  return (
    <>
      <div className='z-10 flex h-20 w-full bg-surface-2' >
        <div className='w-24 md:w-10'>
          <Image src='/title.svg' alt='logo' width={150} height={40} className='h-fix w-fix' />
        </div>
      </div >
    </>
  )
}

export default Header;
