import React from 'react';

type Props = {
    children: React.ReactNode;
    onSubmit?: () => Promise<void>;
};

export const Form: React.FC<Props> = (props) => {
    return (
        <form className="flex flex-col items-center gap-y-4" onSubmit={props.onSubmit ? props.onSubmit : undefined}>
            {props.children}
        </form>
    );
}