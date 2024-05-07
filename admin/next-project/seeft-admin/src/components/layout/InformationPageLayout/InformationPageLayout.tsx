import React from 'react';
import MainLayout from '@components/layout/MainLayout';
import Button from '@components/common/Button';

interface LayoutProps {
  title?: string;
  submitText?: string;
  onClick?: () => void;
  children?: React.ReactNode;
}

export default function InformationPageLayout(props: LayoutProps) {
  return (
    <MainLayout>
      <div className='mx-auto relative md:w-1/2 h-full bg-white-0 p-8'>
        <div className='mx-auto w-fit text-xl text-emphasis mb-8'>
          {props.title}
        </div>
        <div className='flex flex-col gap-3'>
          <div className='my-4 flex flex-col items-center justify-items-end gap-5 text-base text-emphasis'>
            {props.children}
          </div>
        </div>
        <div className='mx-auto w-fit text-emphasis mb-8'>
          <Button className='bg-surface-2 hover:bg-surface-1'
            onClick={props.onClick}>
            {props.submitText}
          </Button>
        </div>
      </div>
    </MainLayout>
  );
}