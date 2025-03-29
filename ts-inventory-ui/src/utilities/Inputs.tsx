interface TextInputProps {
    label: string,
    onChange: (arg0: string)=>void,
    value?: string
}

function TextInput({ label, onChange, value }: TextInputProps) {
    return(
        <div className="relative mt-2">
            <input
                id="text-input"
                className="block px-2.5 pb-2.5 pt-4 w-full text-sm text-l_fg dark:text-fg bg-transparent rounded-lg border-2 border-l_bg2 dark:border-bg2 appearance-none focus:outline-none focus:ring-0 focus:border-l_bg2 dark:focus:border-bg2 peer"
                type="text"
                placeholder=""
                value={value}
                onChange={e => onChange(e.target.value)}
            />
            <label
                htmlFor="text-input"
                className="absolute text-sm text-l_fg dark:text-fg duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] bg-l_bg1 dark:bg-bg1 px-2 peer-focus:px-2 peer-focus:text-l_fg dark:peer-focus:text-fg peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1">
                {label}
            </label>
        </div>
    )
}

interface NumberInputProps {
    label: string,
    onChange: (arg0: number)=>void,
    value?: number
}

function NumberInput({ label, onChange, value }: NumberInputProps) {
    return(
        <div className="relative mt-2">
            <input
               id="number-input"
                className="block px-2.5 pb-2.5 pt-4 w-full text-sm text-l_fg dark:text-fg bg-transparent rounded-lg border-2 border-l_bg2 dark:border-bg2 appearance-none focus:outline-none focus:ring-0 focus:border-l_bg2 dark:focus:border-bg2 peer"
                type="number"
                placeholder=""
                value={value}
                onChange={e => onChange(Number(e.target.value))}
            />
            <label
                htmlFor="number-input"
                className="absolute text-sm text-l_fg dark:text-fg duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] bg-l_bg1 dark:bg-bg1 px-2 peer-focus:px-2 peer-focus:text-l_fg dark:peer-focus:text-fg peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1"
                >
                {label}
            </label>
        </div>
    )
}

export { TextInput, NumberInput }
