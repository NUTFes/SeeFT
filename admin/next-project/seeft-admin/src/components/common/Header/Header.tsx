import Image from 'next/image';
import Router from 'next/router';
import { HeaderProps } from './Header.type';
import { authAtom, userAtom } from '@/store/atoms';
import { useRecoilState } from 'recoil';
import { del } from '@api/signOut';
import { User } from '@type/common';
import { FaUserCircle, FaAngleDown, FaAngleUp } from "react-icons/fa";
import { useEffect, useState, useRef } from 'react';

const Header = (props: HeaderProps) => {
  const [auth, setAuth] = useRecoilState(authAtom);
  const [user, setUser] = useRecoilState(userAtom);
  const [isOpen, setIsOpen] = useState(false);
  const [loadedUserName, setLoadedUserName] = useState('Loading...');
  const userDivRef = useRef<HTMLDivElement>(null);
  const [dropdownWidth, setDropdownWidth] = useState(0);

  useEffect(() => {
    if (user.name) {
      setLoadedUserName(user.name);
    }
  }, [user.name]);

  useEffect(() => {
    if (userDivRef.current) {
      setDropdownWidth(userDivRef.current.offsetWidth);
    }
  }, [userDivRef.current, loadedUserName]);

  const signOut = async () => {
    const signOutUrl: string = process.env.CSR_API_URI + '/mail_auth/web_signout';
    const req = await del(signOutUrl, auth.accessToken);
    const authData = {
      isSignIn: false,
      accessToken: '',
    };
    if (req.status === 200) {
      setAuth(authData);
      setUser({} as User);
      Router.push('/');
    }
  };

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  return (
    <>
      <div className='fixed z-10 flex h-12 w-full flex-row items-center bg-surface-2 border-b-2 border-accent-1 px-2.5'>
        <div className='w-28'>
          <Image src='/title.svg' alt='logo' width={150} height={40} className='h-fit w-fit' />
        </div>
        <div className='ml-auto flex flex-row items-center gap-5 text-lg text-emphasis relative '>
          <div ref={userDivRef} onClick={toggleDropdown} className='flex items-center p-1 gap-2 bg-surface-2 border border-accent-1 rounded-md cursor-pointer whitespace-nowrap hover:bg-surface-1'>
            <FaUserCircle />
            {loadedUserName}
            {isOpen ? <FaAngleDown size={'20px'} /> : <FaAngleUp size={'20px'} />}
          </div>
          {isOpen && (
            <div className='absolute top-full bg-surface-2 border border-accent-1 p-2 shadow-lg rounded-md whitespace-nowrap hover:bg-surface-1' style={{ width: dropdownWidth }}>
              <div onClick={async () => { signOut(); }} className='cursor-pointer'>
                ログアウト
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  )
}

export default Header;
