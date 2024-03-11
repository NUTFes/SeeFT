import React from "react";
import { JsxElement } from "typescript";

interface Props {
    text: string;
    onClick?: () => void;
}

function FailerButton(props: Props): JSX.Element {
    return(
        <>
            <button className='py-3 px-6 rounded-xl text-white-0 bg-caution' onClick={props.onClick}>
                {props.text}
            </button>
        </>
    );
}

export default FailerButton