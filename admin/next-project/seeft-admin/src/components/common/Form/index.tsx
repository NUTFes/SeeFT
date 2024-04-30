import React from 'react';

type Props = {
    children: React.ReactNode;
};

export const Form: React.FC<Props> = (props) => {
    return (
        <form className="flex flex-col items-center gap-y-4">
                {props.children}
        </form>
    );
}