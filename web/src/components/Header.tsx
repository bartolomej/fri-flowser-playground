
import {
    Drawer,
    DrawerContent,
    DrawerTrigger,
} from "@/components/ui/drawer"
import { getData, getFiles } from "../lib/data";
import React, { useState } from "react";
import { Skeleton } from "@/components/ui/skeleton"

const Header = ({ setEditorContent }) => {
    const [fileName, setFileName] = useState<string[]>([]);
    React.useEffect(() => {
        const fetchData = async () => {
            try {
                const files = await getData();
                const names = files.map(file => file.Path.split("/")[1].split(".")[0]);
                console.log(names);
                setFileName(names);
            } catch (error) {
                console.error("Error fetching data:", error);
            }
        };

        fetchData();
    }, []);
    return (
        <header style={styles.header}>
            <Drawer direction="left">
                <DrawerTrigger>
                    <p>Open</p>
                </DrawerTrigger>
                <div className="dark bg-[#213547]">
                    <DrawerContent className="bg-[#213547]" >
                        <div>
                            {
                                fileName.length > 0 ? fileName.map((el: string) => (
                                    <div key={el} className="w-[100%] h-[5vh] hover:cursor-pointer" onClick={()=>setEditorContent(getFiles[el])}>
                                        <p className="text-zinc-400 text-xl">{el}</p>
                                    </div>
                                )) :
                                    <div className="h-[100%] flex flex-col justify-evenly">
                                        <Skeleton className="w-[90%] h-[5vh] rounded-full my-2" />
                                        <Skeleton className="w-[90%] h-[5vh] rounded-full my-2" />
                                        <Skeleton className="w-[90%] h-[5vh] rounded-full my-2" />
                                        <Skeleton className="w-[90%] h-[5vh] rounded-full my-2" />
                                    </div>
                            }
                        </div>
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
        height: "fit-content",
    },
    button: {
        padding: '10px 20px',
        fontSize: '16px',
        cursor: 'pointer',
    },
};

export default Header;