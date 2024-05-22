import ReactJson, {ReactJsonViewProps} from '@microlink/react-json-view'

export function JsonView(props: ReactJsonViewProps) {
    return <ReactJson theme="tomorrow" collapsed={1} {...props} />
}
