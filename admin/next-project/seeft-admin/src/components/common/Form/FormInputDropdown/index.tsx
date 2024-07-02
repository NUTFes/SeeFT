import React from "react";

interface Props {
  text?: string;
  name?: string;
  type?: string;
  onChange?: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  value?: string;
  children?: React.ReactNode;
}

export const FormInputDropdown = (props: Props) => {
  return (
    <div className="flex w-full h-10 items-center">
      <p className="w-40">{props.text}</p>
      <select
        name={props.name}
        onChange={props.onChange}
        value={props.value}
        className="flex-grow h-full bg-transparent border-b border-solid border-accent-1"
      >
        {props.children}
      </select>
    </div>
  );
};
