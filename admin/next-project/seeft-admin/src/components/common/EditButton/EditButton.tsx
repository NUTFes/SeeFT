import clsx from 'clsx';
import React from 'react';
import { MdEdit } from 'react-icons/md';

interface Props {
  className?: string;
  onClick?: () => void;
  children?: React.ReactNode;
}

function EditButton(props: Props): JSX.Element {
  const className = (props.className ? ` ${props.className}` : '');
  return (
    <div className={clsx(className)}>
      <div className='p-2 rounded-full hover:bg-accent-1' onClick={props.onClick}>
        <div className='flex justify-items-center gap-4'>
          <MdEdit />
          <p className='text-center text-sm text-emphasis'>
            {props.children}
          </p>
        </div>
      </div>
    </div>
  );
}

export default EditButton;
