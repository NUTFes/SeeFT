import clsx from 'clsx';
import Head from 'next/head';
import React from 'react';

import { Header } from '@components/common';
import { SideNav } from '@components/common';


interface LayoutProps {
  children?: React.ReactNode;
}

export default function MainLayout(props: LayoutProps) {
  return (
    <>
      <Head>
        <title>SeeFT</title>
        <meta name='' content='' />
        <link rel='icon' href='/icon.ico' />
      </Head>
      <div className={clsx('h-screen w-full')}>
        <div className={clsx('h-12 w-full')}>
          <Header />
        </div>
        <div className='relative overflow-hidden text-emphasis'>
          <div className='transition-all'>
            <SideNav />
          </div>
          <div className={clsx('h-full w-full pl-32 float-right')}>
            {props.children}
          </div>
        </div>
      </div>
    </>
  );
}