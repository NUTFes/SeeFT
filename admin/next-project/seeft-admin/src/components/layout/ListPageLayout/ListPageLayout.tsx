import React from 'react';
import MainLayout from '@components/layout/MainLayout';

interface LayoutProps {
  title?: string;
  children?: React.ReactNode;
}

export default function ListPageLayout(props: LayoutProps) {
  return (
    <MainLayout>
      <div className='w-full h-full bg-white-0 flex-col p-8'>
        <div className='items-center text-xl text-emphasis'>
          {props.title}
        </div>
        <div>
          {props.children}
        </div>
      </div>
    </MainLayout>
  );
}