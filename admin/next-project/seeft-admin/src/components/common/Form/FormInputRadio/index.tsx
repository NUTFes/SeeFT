interface Option {
    value: number | string;
    label: string;
  }
  
  interface Props {
    text?: string;
    name?: string;
    onChange?: () => void;
    value?: number | string;
    options: Option[];
  }

  export const FormInputRadio = (props: Props) => {
    return (
      <div className="flex w-full h-10 items-center">
        <p className="w-40">{props.text}</p>
        <div className="flex-grow h-ful">
          {props.options.map((option) => (
            <label key={option.value} className="inline-flex items-center mr-4">
              <input
                type="radio"
                name={props.name}
                value={option.value}
                checked={props.value === option.value}
                onChange={props.onChange}
                className="bg-transparent border-b border-solid border-accent-1"
              />
              {option.label}
            </label>
          ))}
        </div>
      </div>
    );
  };