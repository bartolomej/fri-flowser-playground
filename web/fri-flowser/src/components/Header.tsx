
import { Drawer, DrawerContent, DrawerTrigger } from "./ui/drawer";

const Header = () => {
    const onButtonClick = () => { }

    return (
        <header style={styles.header}>
            <Drawer direction="left">
                <DrawerTrigger>
                    <button onClick={onButtonClick}>Open</button>
                </DrawerTrigger>
                <div className="dark bg-[#213547]">
                    <DrawerContent className="bg-[#213547]" >
                        
                    </DrawerContent>
                </div>
            </Drawer>
        </header>
    )
}

const styles = {
    header: {
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        padding: '10px 20px',
        backgroundColor: '#282c34',
        color: 'white',
        position: "fixed",
        top: 0,
        left: 0,
       width: "100vw",
    },
    button: {
        padding: '10px 20px',
        fontSize: '16px',
        cursor: 'pointer',
    },
};

export default Header;