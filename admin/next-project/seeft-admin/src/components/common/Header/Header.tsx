import Image from 'next/image';
import Router from 'next/router';
import { HeaderProps } from './Header.type';

const Header = (props: HeaderProps) => {
  return (
    <>
      <div className='fixed z-10 flex h-12 w-full flex-row items-center bg-surface-2 border-b-2 border-accent-1 px-2.5' >
        <div className='w-28'>
          <Image src='/title.svg' alt='logo' width={150} height={40} className='h-fit w-fit' />
        </div>
      </div >
    </>
  )
}

export default Header;
