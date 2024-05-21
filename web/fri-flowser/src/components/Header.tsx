
import { Drawer, DrawerContent, DrawerTrigger } from "./ui/drawer";

const Header = () => {

    return (
        <header style={styles.header}>
            <Drawer direction="left">
                <DrawerTrigger>
                    <p >Open</p>
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
        position: "relative",
        top: 0,
        left: 0,
        width: "100vw",
        height:"fit-content",
    },
    button: {
        padding: '10px 20px',
        fontSize: '16px',
        cursor: 'pointer',
    },
};

export default Header;