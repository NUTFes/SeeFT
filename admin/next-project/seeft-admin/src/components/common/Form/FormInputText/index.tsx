
interface Props {
  text?: string;
  type?: string;
  placeholder?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  value?: string;
}

export const FormInputText = (props: Props) => {
  return (
    <div className="flex w-full h-10 items-center">
      <p className="w-40">{props.text}</p>
      <input
        type={props.type}
        placeholder={props.placeholder}
        className="flex-grow h-full bg-transparent border-b border-solid border-accent-1 text-emphasis"
        onChange={props.onChange}
        value={props.value}
      />
    </div>
  );
}