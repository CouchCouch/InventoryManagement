function TextInput({ label, onChange, value }) {
    return(
        <div className='relative'>
            <input
                id='text-input'
               className='block px-2.5 pb-2.5 pt-4 w-full text-sm text-gray-900 bg-transparent rounded-lg border-1 border-gray-500 appearance-none focus:outline-none peer'
                type='text'
                placeholder=''
                value={value}
                onChange={e => onChange(e.target.value)}
            />
            <label
                htmlFor='text-input'
                className='absolute text-sm text-gray-500 duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] px-2 peer-focus:px-2 peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1'
            >
                {label}
            </label>
        </div>
    )
}

function NumberInput({ label, onChange, value }) {
    return(
        <div className='relative'>
                <input
                id='number-input'
               className='block px-2.5 pb-2.5 pt-4 w-full text-sm text-gray-900 bg-transparent rounded-lg border-1 border-gray-500 appearance-none focus:outline-none peer'
                type='number'
                placeholder=''
                value={value}
                onChange={e => {onChange(e.target.value); console.log(e.target.value)}}
            />
            <label
                htmlFor='number-input'
                className='absolute text-sm text-gray-500 duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] px-2 peer-focus:px-2 peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1'
            >
                {label}
            </label>
        </div>
    )
}

export { TextInput, NumberInput }