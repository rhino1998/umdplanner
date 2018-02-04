import Section from './Section';
import GenEd from './GenEd';

export class Class{
    title: string;
    code: string;
    credits: number;
    prereqs: Array<Class>;

    description: string;
    prerequisite: string;
    restriction: string;

    gened: GenEd;
    sections: Array<Section>;

    conflicts(o: Class):boolean{
        for (let section of this.sections){
            for (let otherSection of o.sections){
                if (!section.conflicts(otherSection)){
                    return false;
                }
            }
        }
        return true;
    }
}

export default Class;

