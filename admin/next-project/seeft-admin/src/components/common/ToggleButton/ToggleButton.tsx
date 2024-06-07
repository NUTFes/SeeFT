import clsx from 'clsx';
import React from 'react';

import s from './Select.module.css';

interface Props {
  className?: string;
  isToggleState: boolean;
  initialLabel?: string;
  toggledLabel?: string;
  onClick?: () => void;
  children?: React.ReactNode;
}

function ToggleButton(props: Props): JSX.Element {
  const className =
    'rounded-full border border-accent-1 py-2 px-4 w-full' +
    (props.isToggleState ? ' bg-surface-1' : ' bg-surface-2') +
    (props.className ? ` ${props.className}` : '');
  return (
    <button className={clsx(className)} onClick={props.onClick}>
      {props.isToggleState ? props.toggledLabel : props.initialLabel}
      {props.children}
    </button>
  );
}

export default ToggleButton;
