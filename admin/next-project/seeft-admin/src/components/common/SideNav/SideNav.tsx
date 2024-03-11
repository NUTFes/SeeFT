import clsx from 'clsx';
import { useRouter } from 'next/router';
import React, { ReactNode, useReducer } from 'react';

import { SeeftLinkItems } from '@constants/linkItem';
import { NavItemProps } from './SideNav.type';

const NavItem = (props: NavItemProps) => {
  let className = '';

  if (props.currentPath === props.href) {
    className = 'flex items-center w-full space-x-2 bg-surface-1 px-2 py-2';
  } else {
    className = 'flex items-center w-full space-x-2 bg-surface-2 px-2 py-2 hover:bg-surface-1/[0.5]';
  }

  if (props.isShow && !props.isParent) {
    className += ' hidden';
  }

  return (
    <a href={props.href}>
      <div className={clsx(className)}>
        <div>{props.icon}</div>
        <div>{props.children}</div>
      </div>
    </a>
  )
}

export default function SideNavigationBar() {
  const router = useRouter();
  const [isSeeftItemsShow, setSeeftItemsShow] = useReducer((state) => !state, false);

  return (
    <div className='fixed right-0 z-10 h-full w-32 bg-surface-2 md:left-0'>
      <div className='text-emphasis'>
        {SeeftLinkItems.map((link) => (
          <div className='border-b-2 border-b-accent-1'>
            <NavItem
              key={link.name}
              icon={link.icon}
              href={link.href}
              currentPath={router.pathname}
              isShow={isSeeftItemsShow}
            >
              {link.name}
            </NavItem>
          </div>
        ))}
      </div>
    </div>
  )
}
