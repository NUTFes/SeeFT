import clsx from 'clsx';
import React from 'react';

import s from './Select.module.css';

interface Props {
  className?: string;
  placeholder?: string;
  value?: string | number;
  defaultValue?: string | number;
  onChange?: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  children?: React.ReactNode;
}

function Select(props: Props): JSX.Element {
  const className =
    'rounded-full border border-accent-2 py-1 px-4 w-full truncate' +
    (props.className ? ` ${props.className}` : '');
  return (
    <div className={clsx(s.customSelect)}>
      <select
        placeholder={props.placeholder}
        className={clsx(s.select, className)}
        value={props.value}
        defaultValue={props.defaultValue}
        onChange={props.onChange}
      >
        {props.children}
      </select>
    </div>
  );
}

export default Select;
