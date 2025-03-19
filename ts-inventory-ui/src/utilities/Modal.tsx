import { XMarkIcon } from "@heroicons/react/24/solid";
import { ReactNode, useEffect } from "react";

interface ModalProps {
    open: boolean,
    onClose: ()=> void,
    title: string,
    children: ReactNode
}

export default function Modal({ open, onClose, title, children }: ModalProps) {
    const escHandler = (e: KeyboardEvent): void => {
        if(e.key === "Escape") {
            onClose()
        }
    }

    useEffect(() => {
        document.addEventListener("keydown", escHandler, false)
    })

    return(
        <div
            onClick={onClose}
            className={
                `fixed inset-0 flex justify-center items-center transition-colors ${open ? "visible bg-bg_dim/40" : "invisible"}
            `}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                className={`
                    bg-bg1 rounded-xl shadow p-6 transition-all w-auto
                    ${open ? "scale-100 opacity-100" : "scale-125 opacity-0"}
                `}
            >
                <div className="block text-center">
                    <h3 className="text-lg font-black text-fg mt-5 mb-1">
                        {title}
                    </h3>
                    <button
                        onClick={onClose}
                        className="btn btn-danger absolute top-1 right-1 p-1 rounded-lg"
                    >
                        <XMarkIcon className="size-6" />
                    </button>
                </div>
                <div className="place-content-center text-center">
                    {children}
                </div>
            </div>
        </div>
    )
}
