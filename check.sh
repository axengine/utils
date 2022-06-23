 for dir in `ls .`
 do
   if [ -d $dir ]
   then
     echo $dir
     cd $dir
     go build
     cd ..
   fi
done

