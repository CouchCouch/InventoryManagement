import PropTypes from 'prop-types'

function TextInput({ label, onChange, value }) {
    return(
        <div className='relative mt-2'>
            <input
                id='text-input'
                className='block px-2.5 pb-2.5 pt-4 w-full text-sm text-gray-900 bg-transparent rounded-lg border-2 border-gray-500 appearance-none focus:outline-none focus:ring-0 focus:border-dark_green peer'
                type='text'
                placeholder=''
                value={value}
                onChange={e => onChange(e.target.value)}
            />
            <label
                htmlFor='text-input'
                className='absolute text-sm text-gray-500 duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] bg-lavender px-2 peer-focus:px-2 peer-focus:text-dark_green peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1'
            >
                {label}
            </label>
        </div>
    )
}

TextInput.propTypes = {
    label: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    value: PropTypes.string
}

function NumberInput({ label, onChange, value }) {
    return(
        <div className='relative mt-2'>
                <input
                id='number-input'
               className='block px-2.5 pb-2.5 pt-4 w-full text-sm text-gray-900 bg-transparent rounded-lg border-2 border-gray-500 appearance-none focus:outline-none focus:ring-0 focus:border-dark_green peer'
                type='number'
                placeholder=''
                value={value}
                onChange={e => {onChange(e.target.value)}}
            />
            <label
                htmlFor='number-input'
                className='absolute text-sm text-gray-500 duration-300 transform -translate-y-4 scale-75 top-2 z-10 origin-[0] bg-lavender px-2 peer-focus:px-2 peer-focus:text-dark_green peer-placeholder-shown:scale-100 peer-placeholder-shown:-translate-y-1/2 peer-placeholder-shown:top-1/2 peer-focus:top-2 peer-focus:scale-75 peer-focus:-translate-y-4 rtl:peer-focus:translate-x-1/4 rtl:peer-focus:left-auto start-1'
            >
                {label}
            </label>
        </div>
    )
}

NumberInput.propTypes = {
    label: PropTypes.string.isRequired,
    onChange: PropTypes.func.isRequired,
    value: PropTypes.oneOfType([
        PropTypes.number,
        PropTypes.string
    ])
}

export { TextInput, NumberInput }
